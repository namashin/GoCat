package main

import (
	"GoCat/app"
	"GoCat/config"
	"GoCat/utils"
	"github.com/getlantern/systray"
	_ "image/jpeg"
)

var appContext *app.RunAnimal

func init() {
	appContext = app.NewRunAnimal(config.Config.RunAnimal, setUpIcons())
}

func main() {
	systray.Run(onReady, onExit)
}

func setUpIcons() map[string][][]byte {
	icons := make(map[string][][]byte)

	animals := []struct {
		name string
		path string
	}{
		{name: "white_cat", path: "./res/cat/white/"},
		{name: "black_cat", path: "./res/cat/black/"},
		{name: "white_horse", path: "./res/horse/white/"},
		{name: "black_horse", path: "./res/horse/black/"},
		{name: "white_parrot", path: "./res/parrot/white/"},
		{name: "black_parrot", path: "./res/parrot/black/"},
	}

	for _, animal := range animals {
		fileCount, err := utils.CountPNGFile(animal.path)
		if err != nil || fileCount <= 0 {
			continue
		}

		icons[animal.name] = utils.LoadIcons(animal.path, animal.name, fileCount)
	}

	return icons
}

func onReady() {
	appContext.Start()
}

func onExit() {
	appContext.End()
}
