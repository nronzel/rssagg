package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nronzel/rssagg/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	// add user to DB
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
