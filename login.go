package main

import (
	"net/http"
	"encoding/json"
	"time"
	"log"

	"github.com/phucfix/chirpy/internal/auth"
	"github.com/phucfix/chirpy/internal/database"
)

const (
	tokenExpirationTime 	   = 1 * time.Hour
	refreshTokenExpirationTime = 60 * 24 * time.Hour
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Password string	`json:"password"`
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

	// Create token
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, tokenExpirationTime)
	if err != nil {
		log.Printf("Can't create json web token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Can't create json web token")
		return
	}

	// Create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Can't create the refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Can't create refresh token")
		return
	}

	dbToken, err := cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpirationTime),
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Can't create token in database")
		respondWithError(w, http.StatusInternalServerError, "Can't create token in database")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: dbToken.Token,
		IsChirpyRed: user.IsChirpyRed,
	})
}
