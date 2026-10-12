package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pippps/httpServer_01/internal/auth"
	"github.com/pippps/httpServer_01/internal/database"
)


func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request){
	chirpID := r.PathValue("chirpID")
	chirp, err := cfg.dbQueries.GetChirp(r.Context(),uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, 404, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil{
		log.Print(err)
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	chirp := &Chirp{
		UserID: userID,
	}
	err = decoder.Decode(chirp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	chirpParams := database.CreateChirpParams{
		Body: chirp.Body,
		UserID: uuid.NullUUID{
			UUID:  chirp.UserID,
			Valid: true,
		},
	}
	chirpResp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 201, chirpResp)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return 
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return 
	}
	chirpID := r.PathValue("chirpID")
	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}
	if dbChirp.UserID.UUID != userID {
		respondWithError(w, 403, "Unauthorized", fmt.Errorf("client asked to delete someone else chirp"))
		return
	}
	if err = cfg.dbQueries.DeleteChirp(r.Context(), uuid.MustParse(chirpID)); err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	w.WriteHeader(204)
}