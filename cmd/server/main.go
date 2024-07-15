package main

import (
	"github.com/khinshankhan/listo/internal/config"
	"github.com/khinshankhan/listo/internal/controller"
	"log"
)

// Version and BuildData get replaced during build with the commit hash and time of build
var (
	CommitHash = ""
	BuildDate  = ""
)

// TODO: build this out to read in a full configuration
func loadConfiguration() *config.Config {
	loadedCfg, err := config.Load(CommitHash, BuildDate)
	if err != nil {
		log.Fatalf("[setup][LoadConfig] %s\n", err.Error())
	}

	return loadedCfg
}

func main() {
	loadedCfg := loadConfiguration()

	controller.Handle(loadedCfg)
}
