package fs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
)

func ListFiles(dir, glob string) []string {
	root := os.DirFS(dir)

	imageFiles, err := fs.Glob(root, glob)

	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, v := range imageFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}

func ReadImageFileToBase64(filePath string) (string, error) {
	// Load and encode the image
	imageData, err := os.ReadFile(filePath)
	// imageData, err := os.ReadFile("data/meter_pic1.png")

	if err != nil {
		errString := fmt.Sprintf("Error reading image file: %v\n", err)
		return "", errors.New(errString)
	}

	// Encode the image as base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return base64Image, nil

}
