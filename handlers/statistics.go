package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type Statistics struct {
	MostPopularFormat     string         `json:"most_popular_format"`
	TopCameraModels       []string       `json:"top_camera_models"`
	UploadFrequencyPerDay map[string]int `json:"upload_frequency_per_day"`
}

func GetStatistics(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secretToken := r.Header.Get("Authorization")
		if secretToken != "your-secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		statistics := Statistics{}

		// Get the most popular image format
		err := db.Get(&statistics.MostPopularFormat, `
            SELECT format FROM images
            GROUP BY format
            ORDER BY COUNT(*) DESC
            LIMIT 1
        `)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the top 10 most popular camera models
		err = db.Select(&statistics.TopCameraModels, `
            SELECT camera_model FROM images
            GROUP BY camera_model
            ORDER BY COUNT(*) DESC
            LIMIT 10
        `)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get image upload frequency per day for the past 30 days
		rows, err := db.Queryx(`
            SELECT upload_date, upload_count FROM upload_statistics
            WHERE upload_date >= $1
        `, time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		uploadFrequency := make(map[string]int)
		for rows.Next() {
			var date string
			var count int
			if err := rows.Scan(&date, &count); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			uploadFrequency[date] = count
		}

		statistics.UploadFrequencyPerDay = uploadFrequency

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statistics)
	}
}
