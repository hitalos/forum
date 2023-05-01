package config

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

var (
	Port  int
	DBURL string
)

func Load() error {

	type Config struct {
		Port  int    `ini:"port" cfg:"port" cfgDefault:"8080" cfgHelper:"Port"`
		DBURL string `ini:"dburl" cfg:"dburl" cfgHelper:"Database URL" cfgRequired:"true"`
	}

	var cfg Config

	config.PrefixEnv = "FORUM"
	config.File = "forum.ini"
	err := config.Parse(&cfg)
	if err != nil {
		return err
	}

	Port = cfg.Port
	DBURL = cfg.DBURL

	return nil
}
