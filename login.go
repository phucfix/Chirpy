package main

import (
	"net/http"
	"encoding/json"
	"time"
	"log"

	"github.com/phucfix/chirpy/internal/auth"
)

const (
	defaultExpiration = 1 * time.Hour
	maxExpiration	  = 1 * time.Hour
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Password string	`json:"password"`
		Email 	 string `json:"email"`
		ExpireIn *int	`json:"expires_in_seconds"`
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

	// Config expire time for jwt
	var expireTime time.Duration

	// If client specify expiration time
	if request.ExpireIn != nil {
		// Invalid expiration
		if *request.ExpireIn <= 0 {
			log.Printf("Invalid expiration time")
			respondWithError(w, http.StatusBadRequest, "Invalid expiration time, must be greate than 0")
			return
		}

		expireTime = time.Duration(*request.ExpireIn) * time.Second
		
		if expireTime > maxExpiration {
			expireTime = defaultExpiration
		}

	} else {
		expireTime = defaultExpiration
	}

	// Create token
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expireTime)
	if err != nil {
		log.Printf("Can't create json web token: %v", err)
		respondWithError(w, http.StatusBadRequest, "Can't create json web token")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})
}
