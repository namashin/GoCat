package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

func CountFiles(dirPath string) (int, error) {
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

func LoadIcons(basePath string, animal string, fileCount int) [][]byte {
	var icons [][]byte

	for i := 0; i < fileCount; i++ {
		icon, err := loadImage(basePath + fmt.Sprintf("%s_%d.png", animal, i))
		if err != nil {
			continue
		}

		icons = append(icons, icon)
	}

	return icons
}

func loadImage(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
