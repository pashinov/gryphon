package handler

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"

	"gryphon/authorizer/config"
	"gryphon/authorizer/models"
)

type Handler struct {
	conf *config.Config
	provider *oidc.Provider
	state string
}

func NewHandler(conf *config.Config) *Handler {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, conf.IdentityServer.Url + conf.IdentityServer.Resources.Identifier)
	if err != nil {
		log.Fatal(err)
	}

	state := tokenGenerator()

	return &Handler{conf, provider, state}
}

func (handler *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	values, ok := r.URL.Query()["client_id"]
	if !ok || len(values) != 1 {
		http.Error(w, "Query param 'client_id' is missing", http.StatusBadRequest)
		return
	}

	clientId := values[0]
	clientSecret, err := getSecret(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(clientSecret)
	log.Println(clientSecret)

	redirectURL := handler.conf.Authorizer.Url + "/oauth/" + clientId + "/callback"

	oauth2Config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     handler.provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL(handler.state), http.StatusFound)
	return
}

func (handler *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := url.Values{}
	data.Set("client_id", r.PostFormValue("client_id"))
	data.Set("grant_type", r.PostFormValue("grant_type"))
	data.Set("refresh_token", r.PostFormValue("refresh_token"))

	clientSecret, err := getSecret(r.PostFormValue("client_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.Set("client_secret", clientSecret)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, handler.conf.IdentityServer.Url + handler.conf.IdentityServer.Resources.RefreshToken, strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode > 299 {
		http.Error(w, "Failed to make request to Identity Server", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload := make(map[string]interface{})
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshTokenResp := models.RefreshTokenResponse{
		AccessToken:      payload["access_token"].(string),
		TokenType:        payload["token_type"].(string),
		RefreshToken:     payload["refresh_token"].(string),
		ExpiresIn:        payload["expires_in"].(float64),
		RefreshExpiresIn: payload["refresh_expires_in"].(float64),
	}

	respData, err := json.MarshalIndent(refreshTokenResp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(respData)

}

func (handler *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	clientId, ok := mux.Vars(r)["client_id"]
	if !ok {
		http.Error(w, "Wrong ClientId", http.StatusBadRequest)
		return
	}

	clientSecret, err := getSecret(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != handler.state {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	redirectURL := handler.conf.Authorizer.Url + "/oauth/" + clientId + "/callback"

	oauth2Config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     handler.provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
	}

	ctx := context.Background()
	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: clientId,
	}
	verifier := handler.provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessTokenResp := models.AccessTokenResponse{
		AccessToken:  oauth2Token.AccessToken,
		TokenType:    oauth2Token.TokenType,
		RefreshToken: oauth2Token.RefreshToken,
		Expiry:       oauth2Token.Expiry,
		UserId:       idToken.Subject,
	}

	data, err := json.MarshalIndent(accessTokenResp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}
