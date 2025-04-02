package main

import (
	"net/http"
	"time"
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"github.com/phucfix/chirpy/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	Body string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID uuid.UUID `json:"user_id"`
}

// Support standard CRUD operations for "chirps". A "chirp" is just a short message that a user can post to the API, like a tweet.

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	var request reqParams
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	// Validate chirpbody
	if !isValidChirpBody(request.Body) {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: profanity(request.Body),
		UserID: request.UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Error creating chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		Body: chirp.Body, 
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Unable to get all chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to get chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
