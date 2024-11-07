package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func setUpIcons() map[string][][]byte {
	icons := make(map[string][][]byte)
	isWindows := runtime.GOOS == "windows"

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
		var fileCount int
		var err error

		if isWindows {
			fileCount, err = countIcoFile(animal.path)
		} else {
			fileCount, err = countPNGFile(animal.path)
		}

		if err != nil || fileCount <= 0 {
			continue
		}

		icons[animal.name] = loadIcons(animal.path, animal.name, fileCount)
	}

	return icons
}

func countPNGFile(dirPath string) (int, error) {
	var fileCount int
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fileCount, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			fileCount++
		}
	}

	return fileCount, nil
}

func countIcoFile(dirPath string) (int, error) {
	var fileCount int
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fileCount, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".ico" {
			fileCount++
		}
	}

	return fileCount, nil
}

func loadIcons(basePath string, animal string, fileCount int) [][]byte {
	var icons [][]byte

	for i := 0; i < fileCount; i++ {
		iconPath := filepath.Join(basePath, fmt.Sprintf("%s_%d.ico", animal, i))
		icon, err := loadIcon(iconPath)
		if err != nil {
			log.Printf("Failed to load icon %s: %v\n", iconPath, err)
			continue
		}
		icons = append(icons, icon)
	}

	return icons
}

func loadIcon(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// getCPUUsage retrieves the current CPU usage percentage.
func getCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return percent[0]
}
