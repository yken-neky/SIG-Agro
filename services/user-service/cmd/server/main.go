package main

import (
	"flag"

	"github.com/sig-agro/services/user-service/server"
)

func main() {
	configFilePath := flag.String("config", "./config/config.json", "Configuration file path")
	flag.Parse()
	server.RunServer(*configFilePath)
}
