package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/snexus/wmeter/entities"
)

type DBConnect struct {
	Db *sql.DB
}

// initDB initializes the SQLite database and creates the todos table if it doesn't exist
func InitDB(dbPath string) *DBConnect {
	db, err := sql.Open("sqlite3", dbPath) // Open a connection to the SQLite database file named app.db
	if err != nil {
		log.Fatal(err) // Log an error and stop the program if the database can't be opened
	}

	// SQL statement to create the todos table if it doesn't exist
	sqlStmt := `
 CREATE TABLE IF NOT EXISTS meterdata (
  hash TEXT NOT NULL PRIMARY KEY,
  path TEXT,
  taken_timestamp TEXT,
  inserted_timestamp TEXT, 
  meter INTEGER
 );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt) // Log an error if table creation fails
	}
	return &DBConnect{Db: db}
}

func (c *DBConnect) InsertRecord(image entities.ImageMetadata) error {
	query := "INSERT INTO meterdata (hash, path, taken_timestamp, inserted_timestamp, meter) VALUES (?,?,?,?, ?)"
	_, err := c.Db.Exec(query, image.Hash, image.ImagePath, image.TakenTimestamp, time.Now(), image.MeterReading)
	if err != nil {
		return err
	}
	return nil

}

func processRows(rows *sql.Rows) ([]entities.ImageMetadata, error) {
	data := []entities.ImageMetadata{}

	for rows.Next() {
		i := entities.ImageMetadata{}
		var dateStrTaken, dateStrInserted string // Temporary string to hold the date

		// Scan into the temporary string for date
		err := rows.Scan(&i.Hash, &i.ImagePath, &dateStrTaken, &dateStrInserted, &i.MeterReading)
		if err != nil {
			fmt.Println("Got error ", err)
			return nil, err
		}

		// Convert string to time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05+08:00", dateStrTaken) // Adjust format as needed
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %s: %v", dateStrInserted, err)
		}
		i.TakenTimestamp = parsedTime

		data = append(data, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil

}

// Query the database for last X days of data
func (c *DBConnect) QueryLatestDays(ndays int) ([]entities.ImageMetadata, error) {
	interval := fmt.Sprintf("-%d days", ndays)

	rows, err := c.Db.Query(
		"SELECT * FROM meterdata WHERE datetime(inserted_timestamp) >= datetime('now', ?) ORDER by taken_timestamp ASC;",
		interval)

	if err != nil {
		return nil, err
	}

	return processRows(rows)
}
