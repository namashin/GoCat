package main

import (
	"github.com/go-ini/ini"
	"log"
)

const ConfigPath = "./config.ini"

type ConfList struct {
	RunAnimal string
}

var Config ConfList

func init() {
	cfg, err := ini.Load(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	Config = ConfList{
		RunAnimal: cfg.Section("animal").Key("run_animal").MustString("white_cat"),
	}
}
