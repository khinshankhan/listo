package main

import (
	"github.com/khinshankhan/listo/internal/config"
	"github.com/khinshankhan/listo/internal/controller"
	"github.com/khinshankhan/listo/internal/service/log"
	"go.uber.org/zap"
)

// Version and BuildData get replaced during build with the commit hash and time of build
var (
	CommitHash = ""
	BuildDate  = ""
)

// TODO: build this out to read in a full configuration
func loadConfiguration() *config.Config {
	logger := log.NewLogger()

	loadedCfg, err := config.Load(CommitHash, BuildDate)
	if err != nil {
		logger.Fatal(
			"Config is broken",
			zap.Error(err),
		)
	}

	return loadedCfg
}

func main() {
	loadedCfg := loadConfiguration()

	controller.Handle(loadedCfg)
}
