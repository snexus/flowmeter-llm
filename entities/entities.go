package entities

import "time"

type MeterReadingResult struct {
	Steps []struct {
		Explanation string `json:"explanation"`
		Output      string `json:"output"`
	} `json:"steps"`
	MeterReading string `json:"meter_reading"`
}

type ImageMetadata struct {
	ImagePath      string    `json:"image_path"`
	TakenTimestamp time.Time `json:"date_taken"`
	MeterReading   int       `json:"meter_reading"`
	Hash           string    `json:"hash"`
}
