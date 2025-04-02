package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerResetHit(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Reset OK"))
}

func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You can not do this")
		return
	}

	if err := cfg.dbQueries.DeleteUsers(req.Context()); err != nil {
		log.Printf("Can not delete user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Can not delete user")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(`Delete all users in the database successfully`))
}
