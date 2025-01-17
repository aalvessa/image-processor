package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aalvessa/image-processor/interfaces/handlers"
	repositories "github.com/aalvessa/image-processor/interfaces/repositores"
	"github.com/aalvessa/image-processor/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func main() {

	db, err := utils.ConnectDB()

	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	defer db.Close()

	// Run migrations
	if err := utils.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	startRouter(db)
}

func startRouter(db *sqlx.DB) {
	linkRepo := repositories.NewLinks(db)
	imageRepo := repositories.NewImages(db)
	statisticsRepo := repositories.NewStatistics(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/upload-link", handlers.GenerateUploadLink(linkRepo))
	r.Post("/upload", handlers.UploadImage(linkRepo, imageRepo, statisticsRepo))
	r.Get("/image/{id}", handlers.GetImage(imageRepo))
	r.Get("/statistics", handlers.GetStatistics(statisticsRepo))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hola")) })

	http.ListenAndServe(":8000", r)
}
