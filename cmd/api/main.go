package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/jumaniyozov/gores/internal/app/api"
	"log"
)

var (
	configPath string = "configs/api.toml"
)

func init() {
	flag.StringVar(
		&configPath,
		"path",
		"configs/api.toml",
		"path to config file in .toml format",
	)
}

func main() {
	// Initialize and parse flag values from terminal api.exe -path configs/<name>.toml
	flag.Parse()
	config := api.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Println("Can not find configs file. Using default values.", err)
	}

	server := api.New(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
