package main

import (
	"time"
	"net/http"
	"log"

	"github.com/phucfix/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Can't get refresh token from request")
		respondWithError(w, http.StatusBadRequest, "Can't get refresh token")
		return
	}

	// Validate the refresh token
	dbRefreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(req.Context(), refreshToken)
	if err != nil {
		log.Printf("Error validating the refresh token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't validate the refresh token")
		return
	}

	if time.Now().After(dbRefreshToken.ExpiresAt) {
		log.Printf("Refresh token is expired")
		respondWithError(w, http.StatusUnauthorized, "Expired refresh token")
		return
	}

	if dbRefreshToken.RevokedAt.Valid {
		log.Printf("Refresh token is revoked")
		respondWithError(w, http.StatusUnauthorized, "Revoked refresh token")
		return	
	}

	// Create new token
	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.jwtSecret, tokenExpirationTime)
	if err != nil {
		log.Printf("Can't create json web token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't create json web token")
		return
	}

	type Response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, Response{Token: token})
}
