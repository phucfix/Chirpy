package main

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebHooks(w http.ResponseWriter, req *http.Request) {
	type Data struct {
		UserId string `json:"user_id"`
	}

	type ReqParams struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	var request ReqParams
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Can't decode json: %v", err)
		respondWithError(w, http.StatusBadRequest, "Can't decode json")
		return
	}

	// If event is not "user.upgraded", return immediately
	if request.Event != "user.upgraded" {
		log.Printf("Event is not user upgrade")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(request.Data.UserId)
	if err != nil {
		log.Printf("Error parsing uuid: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error parsing uuid")
		return
	}

	// Upgrade user in database
	err = cfg.dbQueries.UpgradeUserById(req.Context(), userId)
	if err != nil {
		log.Printf("User is not upgrade successfully: %v", err)
		respondWithError(w, http.StatusInternalServerError, "User is not upgrade successfully")
	}

	w.WriteHeader(http.StatusNoContent)
}
