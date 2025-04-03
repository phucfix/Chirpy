package main

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/phucfix/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Password string `json:"password"`
		Email 	 string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	var request reqParams
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}
	
	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), request.Email)
	if err != nil {
		log.Printf("Error authenticate user: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't recognize user")
		return
	}

	// Check matching password
	if err := auth.CheckPasswordHash(user.HashedPassword, request.Password); err != nil {
		log.Printf("password comparison errors: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
