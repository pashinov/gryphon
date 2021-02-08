package main

import (
	"log"
	"strconv"

	"github.com/jessevdk/go-flags"

	"gryphon/authorizer/config"
	"gryphon/authorizer/server"
)

func main() {
	var opts struct {
		Help   bool   `short:"h" long:"help" description:"show help message"`
		Config string `short:"c" long:"config" default:"config/backend-app.json"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Parsing arguments: ", err.Error())
		return
	}

	conf := config.NewConfig()
	err = conf.Init(opts.Config)
	if err != nil {
		log.Fatal("Failed to initialize of configuration file: ", err)
		return
	}

	serv := server.NewServer().Init(conf)
	serv.Run(conf.Server.Host + ":" + strconv.FormatUint(uint64(conf.Server.Port), 10))
}
