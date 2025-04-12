package usecases

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"

	openai "github.com/sashabaranov/go-openai"
	"github.com/snexus/wmeter/entities"
	"github.com/snexus/wmeter/exif"
	"github.com/snexus/wmeter/fs"
	"github.com/snexus/wmeter/llm"
)

func FetchLatestImageFilesFromDisk(dir, glob string, maxImages int) map[string]entities.ImageMetadata {
	images := []entities.ImageMetadata{}

	imageMap := make(map[string]entities.ImageMetadata, maxImages)

	imageFiles := fs.ListFiles(dir, glob)

	for _, path := range imageFiles {
		time, err := exif.GetExifDateTaken(path)
		if err != nil {
			fmt.Printf("Error getting EXIF date taken: %v\n", err)
			continue
		}
		base64, err := fs.ReadImageFileToBase64(path)

		if err != nil {
			fmt.Printf("Couldn't calculate base64 for the file: %v\n", err)
			continue
		}
		images = append(images, entities.ImageMetadata{ImagePath: path, TakenTimestamp: time, Hash: CreateMd5Hash(base64)})

	}
	slices.SortFunc(images, func(a, b entities.ImageMetadata) int { return -a.TakenTimestamp.Compare(b.TakenTimestamp) })

	nImages := min(len(images), maxImages)

	
	// fmt.Println("Images: ", images[:nImages])

	for _, image := range images[:nImages] {
		imageMap[image.Hash] = image
	}

	return imageMap
	// return images[:nImages]
}

func AnalyzeMultipleImages(client *openai.Client, images []entities.ImageMetadata, prompt string, modelName string) []entities.ImageMetadata {
	for i, image := range images {
		out1, err := llm.DescribeImage(client, image.ImagePath, prompt, modelName)
		if err != nil {
			log.Fatal(err)

		}
		reading, err := strconv.Atoi(out1.MeterReading)

		if err != nil {
			log.Fatal(err)
		}

		images[i].MeterReading = reading
	}
	return images
}

func CreateMd5Hash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// Assumes imagedata is sorted by date taken
func CalculateFlow(images []entities.ImageMetadata, maxImages int) {

	slices.SortFunc(images, func(a, b entities.ImageMetadata) int { return -a.TakenTimestamp.Compare(b.TakenTimestamp) })
	fmt.Printf("\nFlow summary over last %d days:\n", maxImages-1)

	nImages := min(len(images), maxImages)
	dailyConsumption := make([]float64, nImages-1)

	for i := range nImages - 1 {
		diff_hours := images[i].TakenTimestamp.Sub(images[i+1].TakenTimestamp).Hours()
		litres_day := float64(images[i].MeterReading-images[i+1].MeterReading) / diff_hours * 24.0
		fmt.Printf("Flow [L/day] between %s (%d) and %s (%d): %.1f\n", 
			images[i+1].TakenTimestamp.Format("2006-01-02 15:04"), images[i+1].MeterReading, 
			images[i].TakenTimestamp.Format("2006-01-02 15:04"), images[i].MeterReading, 
			litres_day)
		dailyConsumption[i] = litres_day
	}
	average := calculateAverage(dailyConsumption)
	fmt.Printf("\nAverage flow [L/day] over last %d days: %.1f [L/day]\n", maxImages-1, average)

}

func calculateAverage(list []float64) float64 {
	total := 0.0
	for _, num := range list {
		total += num
	}
	return total / float64(len(list))
}
