package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gryphon/authorizer/config"
	"gryphon/authorizer/routes"
)

type Server struct {
	conf   *config.Config
	router *mux.Router
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Init(conf *config.Config) *Server {
	server.conf = conf

	server.router = routes.NewRouter(server.conf)

	return server
}

func (server *Server) Run(addr string) {
	log.Println("Listening to ", addr)
	log.Fatal(http.ListenAndServe(addr, server.router))
}
