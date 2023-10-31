package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not defined in environment variables")
	}

	connstr := os.Getenv("CONNSTR")
	if connstr == "" {
		log.Fatal("CONNSTR not defined in environment variables")
	}

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	cfg := apiConfig{
		DB: dbQueries,
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", cfg.handlerUsersCreate)
	v1Router.Get("/users", cfg.middlewareAuth(cfg.handlerUsersGet))

	v1Router.Post("/feeds", cfg.middlewareAuth(cfg.handlerFeedsCreate))
	v1Router.Get("/feeds", cfg.handlerFeedsGet)

	v1Router.Post("/feed_follows", cfg.middlewareAuth(cfg.handlerFeedFollowsCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handlerFeedFollowsDelete))
	v1Router.Get("/feed_follows", cfg.middlewareAuth(cfg.handlerFeedFollowsGet))

	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerError)

	r.Mount("/v1", v1Router)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startHarvest(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Server started on localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
