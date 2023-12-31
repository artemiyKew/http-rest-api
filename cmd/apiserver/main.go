package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/artemiyKew/http-rest-api/internal/app"
	"github.com/artemiyKew/http-rest-api/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/apiserver.toml", "path to config file")
	app.Migrate()
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
