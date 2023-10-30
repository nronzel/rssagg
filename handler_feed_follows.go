package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nronzel/rssagg/internal/database"
)

func (cfg *apiConfig) handlerFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	feedFollow, err := cfg.createFeedFollow(r.Context(), user.ID, params.FeedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowtoFeedFollow(feedFollow))
}

func (cfg *apiConfig) createFeedFollow(ctx context.Context, userID, feedID uuid.UUID) (database.FeedFollow, error) {
	return cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIDString := chi.URLParam(r, "feedFollowID")
	feedID, err := uuid.Parse(feedIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse feed ID")
		return
	}

	removedFeedFollow, err := cfg.DB.DeleteFeedFollow(r.Context(), feedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowtoFeedFollow(removedFeedFollow))
}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetAllFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get feed follows")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowstoFeedFollows(feedFollows))
}
