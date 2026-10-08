package main

import (
	"net/http"
	"fmt"
)

func (cfg *apiConfig) countHandler(w http.ResponseWriter, r *http.Request) {
	message_string := fmt.Sprintf(
		"<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(message_string))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	if cfg.platform != "dev" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		return
	}
	if err := cfg.dbQueries.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, 400, "error deleting users", err)
	}
	w.Write([]byte("OK"))
}