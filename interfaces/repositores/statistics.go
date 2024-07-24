package repositories

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type statistics struct {
	db *sqlx.DB
}

type Statistics struct {
	MostPopularFormat     *string        `json:"most_popular_format"`
	TopCameraModels       []string       `json:"top_camera_models"`
	UploadFrequencyPerDay map[string]int `json:"upload_frequency_per_day"`
}

func NewStatistics(db *sqlx.DB) *statistics {
	return &statistics{
		db: db,
	}
}

func (s *statistics) UpdateUploadStatistics() {
	currentDate := time.Now().Format("2006-01-02")
	_, err := s.db.Exec(`
        INSERT INTO upload_statistics (upload_date, upload_count)
        VALUES ($1, 1)
        ON CONFLICT (upload_date)
        DO UPDATE SET upload_count = upload_statistics.upload_count + 1
    `, currentDate)
	if err != nil {
		fmt.Printf("error updating upload statistics: %v\n", err)
	}
}

func (s *statistics) GetStatistics() (*Statistics, error) {

	statistics := Statistics{}

	// Get the most popular image format
	err := s.db.Get(&statistics.MostPopularFormat, `
		SELECT format FROM images
		GROUP BY format
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`)
	if err != nil {
		return nil, fmt.Errorf("getting most popular image")
	}

	// Get the top 10 most popular camera models
	err = s.db.Select(&statistics.TopCameraModels, `
		SELECT camera_model FROM images
		GROUP BY camera_model
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("getting 10 most popular camera model")
	}

	// Get image upload frequency per day for the past 30 days
	rows, err := s.db.Queryx(`
		SELECT upload_date, upload_count FROM upload_statistics
		WHERE upload_date >= $1
	`, time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("getting image frequency")
	}
	defer rows.Close()

	uploadFrequency := make(map[string]int)
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("counting frequency")
		}
		uploadFrequency[date] = count
	}

	statistics.UploadFrequencyPerDay = uploadFrequency

	return &statistics, nil
}
