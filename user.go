package main

import (
	"net/http"
	"log"
	"time"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/phucfix/chirpy/internal/auth"
	"github.com/phucfix/chirpy/internal/database"
)

type User struct {
	ID        		uuid.UUID `json:"id"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedAt 		time.Time `json:"updated_at"`
	Email     		string    `json:"email"`
	Token 	  		string 	  `json:"token"`
	RefreshToken	string	  `json:"refresh_token"` 
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	var request reqParams
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	// Create new user in database
	hashPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email: request.Email,
		HashedPassword: hashPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt, 
		Email: user.Email,
	})
}
