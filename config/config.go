package config

import (
	"GoCat/app"
	"github.com/go-ini/ini"
	"os"
)

type ConfList struct {
	RunAnimal string
}

var Config ConfList

func init() {
	cfg, err := ini.Load(app.ConfigPath)
	if err != nil {
		os.Exit(1)
	}

	Config = ConfList{
		RunAnimal: cfg.Section("animal").Key("run_animal").MustString("white_cat"),
	}
}
