package models

import "time"

type AccessTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	UserId       string    `json:"user_id"`
}

type RefreshTokenResponse struct {
	AccessToken      string  `json:"access_token"`
	TokenType        string  `json:"token_type"`
	RefreshToken     string  `json:"refresh_token"`
	ExpiresIn        float64 `json:"expires_in"`
	RefreshExpiresIn float64 `json:"refresh_expires_in"`
}
