package main

import (
	"fmt"
	"github.com/khinshankhan/listo/config"
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

	fmt.Println("Loaded configuration:")
	fmt.Printf("\tBuildDate: %v\n", loadedCfg.Meta.BuildDate)
	fmt.Printf("\tCommitHash: %v\n", loadedCfg.Meta.CommitHash)
}
