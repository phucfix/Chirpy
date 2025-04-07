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
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
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

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Email 	string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	var request reqParams
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	// Validate access token in header
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't validate JWT")
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

	hashPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update user infomation in database
	err = cfg.dbQueries.UpdatePersonalInformation(req.Context(), database.UpdatePersonalInformationParams{
		Email: request.Email,
		HashedPassword: hashPassword,
		ID: userID,
	})
	if err != nil {
		log.Printf("Error editing user: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Error editing user")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), request.Email)
	if err != nil {
		log.Printf("Can't identify the user: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't indentify the user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
