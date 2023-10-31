package main

import (
	"net/http"
	"strconv"

	"github.com/nronzel/rssagg/internal/database"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	const defaultLimit = 10

	limitString := r.URL.Query().Get("limit")
	limit := 10
	if specifiedLimit, err := strconv.Atoi(limitString); err == nil {
		limit = specifiedLimit
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem fetching posts")
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
