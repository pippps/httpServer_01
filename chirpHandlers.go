package main
import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/google/uuid"

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
	decoder := json.NewDecoder(r.Body)
	chirp := &Chirp{}
	err := decoder.Decode(chirp)
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