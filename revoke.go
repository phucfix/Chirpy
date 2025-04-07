package main

import (
	"net/http"
	"log"
	"time"

	"github.com/phucfix/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Can't get refresh token from request")
		respondWithError(w, http.StatusBadRequest, "Can't get refresh token")
		return
	}

	dbRefreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(req.Context(), refreshToken)
	if err != nil {
		log.Printf("Error validating the refresh token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't validate the refresh token")
		return
	}

	// Check if refresh token is expire first
	if time.Now().After(dbRefreshToken.ExpiresAt) {
		log.Printf("Refresh token is expired")
		respondWithError(w, http.StatusUnauthorized, "Expired refresh token")
		return
	}

	err = cfg.dbQueries.RevokeToken(req.Context(), dbRefreshToken.Token)
	if err != nil {
		log.Printf("Can't revoke the token: %v", err)
		respondWithError(w, http.StatusBadRequest, "Can't revoke the token")
		return
	}

	w.WriteHeader(204)
}
