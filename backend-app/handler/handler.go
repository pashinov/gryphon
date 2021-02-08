package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gryphon/backend-app/config"
)

type Handler struct {
	conf *config.Config
}

func NewHandler(conf *config.Config) *Handler {
	return &Handler{conf}
}

func (handler *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{

	}

	req, err := http.NewRequest(http.MethodGet, handler.conf.IdentityServer.Url+handler.conf.IdentityServer.Resources.UserInfo, nil)
	if err != nil {
		http.Error(w, "Failed to create http request to get user info", http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", r.Header.Get("Authorization"))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to make http request to get user info", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	payload := make(map[string]interface{})
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		Name:     payload["name"].(string),
		Username: payload["preferred_username"].(string),
		Email:    payload["email"].(string),
	}

	respData, err := json.MarshalIndent(userInfo, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(respData)
}
