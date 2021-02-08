 package handler

 import (
	 "os"
 )

 var clients map[string]string

 func init() {
	 clients = make(map[string]string)
	 backendAppClientId := os.Getenv("BACKEND_APP_CLIENT_ID")
	 backendAppClientSecret := os.Getenv("BACKEND_APP_CLIENT_SECRET")
	 clients[backendAppClientId] = backendAppClientSecret
 }
