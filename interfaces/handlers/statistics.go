package handlers

import (
	"encoding/json"
	"net/http"

	repositories "github.com/aalvessa/image-processor/interfaces/repositores"
)

type StatisticsGetter interface {
	GetStatistics() (*repositories.Statistics, error)
}

func GetStatistics(statisticsGettter StatisticsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secretToken := r.Header.Get("Authorization")
		if secretToken != "secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		statistics, err := statisticsGettter.GetStatistics()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statistics)
	}
}
