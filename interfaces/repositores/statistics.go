package repositories

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type statistics struct {
	db *sqlx.DB
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
