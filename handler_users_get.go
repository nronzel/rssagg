package main

import (
	"net/http"

	"github.com/nronzel/rssagg/internal/database"
)

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
