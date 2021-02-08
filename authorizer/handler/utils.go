package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
)

func tokenGenerator() string {
	c := 10
	b := make([]byte, c)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getSecret(clientId string) (string, error) {
	if _, ok := clients[clientId]; ok {
		return clients[clientId], nil
	}
	return "", errors.New("Undefined clientId: " + clientId)
}
