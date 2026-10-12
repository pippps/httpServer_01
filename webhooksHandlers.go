package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pippps/httpServer_01/internal/auth"
)

func (cfg *apiConfig) upgradeHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if apiKey != cfg.polkaKey {
		w.WriteHeader(401)
		return 
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := &requestBody{}

	if err = decoder.Decode(reqBody); err != nil {
		w.WriteHeader(400)
		return
	}
	if reqBody.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	if err := cfg.dbQueries.UpgradeChirpyRed(r.Context(), uuid.MustParse(reqBody.Data.UserID)); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
