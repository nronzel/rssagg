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

	apiCfg := apiConfig{
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

	v1Router.Post("/users", apiCfg.handlerUsersCreate)

	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerError)

	r.Mount("/v1", v1Router)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Server started on localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
