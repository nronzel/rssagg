package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("PORT not defined in environment variables")
	}
	port := os.Getenv("PORT")

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	apiRouter := chi.NewRouter()

	r.Mount("/v1", apiRouter)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Server started on localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
