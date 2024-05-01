package main

import (
	"GoCat/app"
	"GoCat/utils"
	"github.com/getlantern/systray"
	_ "image/jpeg"
)

func main() {
	systray.Run(onReady, onExit)
}

var AnimalIcons map[string][][]byte = make(map[string][][]byte)

func init() {
	setUp()
}

func setUp() {
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
		fileCount, err := utils.CountFiles(animal.path)
		if err != nil {
			continue
		}

		AnimalIcons[animal.name] = utils.LoadIcons(animal.path, animal.name, fileCount)
	}
}

func onReady() {
	runAnimal := app.NewRunAnimal("white_cat", AnimalIcons)
	runAnimal.Start()
}

func onExit() {

}
