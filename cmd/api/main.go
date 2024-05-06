package main

import (
	"bash_api/config"
	"bash_api/internal/httpServer"
	"log"
)

// @title           Bash App API
// @version         1.0
// @description     This is a sample bash server.

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	viperInstance, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot load config. Error: {%s}", err.Error())
	}

	cfg, err := config.ParseConfig(viperInstance)
	if err != nil {
		log.Fatalf("Cannot parse config. Error: {%s}", err.Error())
	}

	s := httpServer.NewServer(cfg)
	if err = s.Run(); err != nil {
		log.Print(err)
	}
}
