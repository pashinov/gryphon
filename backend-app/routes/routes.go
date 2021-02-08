package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"gryphon/backend-app/config"
	"gryphon/backend-app/handler"
	"gryphon/backend-app/middleware"
)

func NewRouter(conf *config.Config) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	h := handler.NewHandler(conf)

	r.Use(middleware.Logging)
	r.Use(middleware.TokenAuthMiddleware)

	r.HandleFunc("/backend-app/user/info", h.GetUserInfo).Methods(http.MethodGet, http.MethodOptions)

	return r
}
