package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nronzel/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("PORT not defined in environment variables")
	}
	port := os.Getenv("PORT")
	connstr := os.Getenv("CONNSTR")
	db, err := sql.Open("postgres", connstr)

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	apiRouter := chi.NewRouter()

	apiRouter.Get("/readiness", handlerReadiness)
	apiRouter.Get("/err", handlerError)

    apiRouter.Post("/users", apiCfg.handlerUserCreate)

	r.Mount("/api/v1", apiRouter)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Server started on localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
