package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nronzel/rssagg/internal/database"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create feed")
		return
	}

	// create feed follow
	feedFollow, err := cfg.createFeedFollow(r.Context(), user.ID, feed.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create feed follow")
		return
	}

	response := map[string]interface{}{
		"feed":        databaseFeedToFeed(feed),
		"feed_follow": databaseFeedFollowtoFeedFollow(feedFollow),
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem retrieving feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
