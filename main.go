package main

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

type Config struct {
	DatabaseURL string `ini:"database_url" cfg:"database_url" cfgRequired:"true" cfgHelper:"Database URL"`
}

func main() {
	cfg := Config{}

	config.File = "config.ini"
	err := config.Parse(&cfg)
	if err != nil {
		println(err)
		return
	}

	println(cfg.DatabaseURL)
}
