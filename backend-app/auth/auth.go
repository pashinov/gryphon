package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
)

var verifier *oidc.IDTokenVerifier

func init() {
	clientID := os.Getenv("CLIENT_ID")
	identifierURL := os.Getenv("IDENTIFIER_URL")

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, identifierURL)
	if err != nil {
		panic(err)
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier = provider.Verifier(oidcConfig)
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	ctx := context.Background()
	_, err := verifier.Verify(ctx, tokenString)
	return err
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
