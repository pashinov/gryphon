package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"gryphon/authorizer/config"
	"gryphon/authorizer/handler"
)

func NewRouter(conf *config.Config) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	h := handler.NewHandler(conf)

	r.HandleFunc("/oauth/token", h.GetToken).Methods(http.MethodGet)
	r.HandleFunc("/oauth/token", h.RefreshToken).Methods(http.MethodPost)
	r.HandleFunc("/oauth/{client_id}/callback", h.Callback).Methods(http.MethodGet)

	return r
}

