package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/pippps/httpServer_01/internal/auth"
	"github.com/pippps/httpServer_01/internal/database"
)

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserResp struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	userJson := &UserReq{}
	err := decoder.Decode(userJson)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userJson.Password)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userJson.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, user)
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userJson := &UserReq{}
	err := decoder.Decode(userJson)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	dbUser, err := cfg.dbQueries.GetUser(r.Context(), userJson.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	ok, err := auth.CheckPasswordHash(userJson.Password, dbUser.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	userResp := UserResp{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, 200, userResp)
}
