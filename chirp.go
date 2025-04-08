package main

import (
	"net/http"
	"time"
	"encoding/json"
	"log"
	"strings"
	"sort"

	"github.com/google/uuid"

	"github.com/phucfix/chirpy/internal/database"
	"github.com/phucfix/chirpy/internal/auth"
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
		Body 	string `json:"body"`
		UserID 	uuid.UUID `json:"user_id"`
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

	// Validate JWT
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error get bearer token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't get bearer token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't validate JWT")
		return
	}

	if userID == uuid.Nil {
		log.Printf("Can't identify the user")
		respondWithError(w, http.StatusUnauthorized, "Can't validate user") 
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: profanity(request.Body),
		UserID: userID,
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
	// Optional query parameter named author_id
	// If the author_id query parameter is provided, the endpoint should return only the chirps for that author
	// Else the endpoint should return all chirps as it did before

	// Grab the query parameters from the URL
	authorId := req.URL.Query().Get("author_id")
	
	sortQuery := req.URL.Query().Get("sort")
	// if it exists, or an empty string if it doesn't

	var chirps []database.Chirp
	var err error
	if authorId == "" {
		chirps, err = cfg.dbQueries.GetChirps(req.Context())
		if err != nil {
			log.Printf("Unable to get all chirps: %v", err)
			respondWithError(w, http.StatusBadRequest, "Unable to get chirps")
			return
		}
	} else {
		authorUUID, err := uuid.Parse(authorId)
		if err != nil {
			log.Printf("Can't parse author uuid: %v", err)
			respondWithError(w, http.StatusBadRequest, "Can't parse author uuid")
			return
		}

		chirps, err = cfg.dbQueries.GetChirpsByUserId(req.Context(), authorUUID)
		if err != nil {
			log.Printf("Unable to get chirps by author id: %v", err)
			respondWithError(w, http.StatusNotFound, "Unable to get chirps by author")
			return
		}
	}

	
	var chirpResp []Chirp
	for _, v := range chirps {
		chirpResp = append(chirpResp, Chirp{
			ID: v.ID,
			Body: v.Body,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			UserID: v.UserID,
		})
	}

	// Sort the chirps
	if strings.ToLower(sortQuery) == "desc" {
		
		sort.Slice(chirpResp, func(i, j int) bool {
			return chirpResp[j].CreatedAt.Before(chirpResp[i].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpResp)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Unable to parse ID string to UUID: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to parse to UUID")
		return
	}

	chirp, err := cfg.dbQueries.GetChirpById(req.Context(), chirpID)
	if err != nil {
		log.Printf("Unable to get chirp by ID: %v", err)
		respondWithError(w, http.StatusNotFound, "Unable to get chirp by ID: %v") 
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		Body: chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	// Validate JWT
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error get bearer token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't get bearer token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't validate JWT")
		return
	}

	if userID == uuid.Nil {
		log.Printf("Can't identify the user")
		respondWithError(w, http.StatusUnauthorized, "Cant' validate user")
		return
	}

	// Get chirp
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Unable to parse ID string to UUID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Unable to parse to UUID")
		return
	}

	chirp, err := cfg.dbQueries.GetChirpById(req.Context(), chirpID)
	if err != nil {
		log.Printf("Unable to get chirp by ID: %v", err)
		respondWithError(w, http.StatusNotFound, "Unable to get chipr by ID")
		return
	}
	
	if userID != chirp.UserID {
		log.Printf("Not allowed to do that")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete the chirp
	err = cfg.dbQueries.DeleteChirpById(req.Context(), database.DeleteChirpByIdParams{
		ID: chirp.ID,
		UserID: chirp.UserID,
	})	
	if err != nil {
		log.Printf("Unable to delete chirp: %v", err)
		respondWithError(w, http.StatusForbidden, "Unable to delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
