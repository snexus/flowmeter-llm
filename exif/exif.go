package exif

import (
	// "fmt"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"os"
	"time"
)

func GetExifDateTaken(filePath string) (time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return time.Time{}, err
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
		return time.Time{}, err
	}
	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime()
	// fmt.Println("Taken: ", tm)

	// lat, long, _ := x.LatLong()
	// fmt.Println("lat, long: ", lat, ", ", long)
	return tm, nil
}
