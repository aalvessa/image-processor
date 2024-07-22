package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aalvessa/image-processor/handlers"
	"github.com/aalvessa/image-processor/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/generate-upload-link", handlers.GenerateUploadLink(db))
	r.Post("/upload", handlers.UploadImage(db))
	r.Get("/image/{id}", handlers.GetImage(db))
	r.Get("/statistics", handlers.GetStatistics(db))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hola")) })

	http.ListenAndServe(":8000", r)
}
