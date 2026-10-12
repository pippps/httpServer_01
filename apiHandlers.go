package main

import (
	"encoding/json"
	"fmt"
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
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	Token        string
	RefreshToken string
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 402, "Something went wrong", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, 402, "Something went wrong", err)
		return
	}

	userResp := UserResp{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, 200, userResp)
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type tokenResponse struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "something went wrong", err)
		return
	}
	dbRT, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "something went wrong", err)
		return
	}
	if dbRT.RevokedAt.Valid {
		respondWithError(w, 401, "unauthorized", fmt.Errorf("token revoked"))
		return
	}
	token, err := auth.MakeJWT(dbRT.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "something went wrong", err)
		return
	}
	respondWithJSON(w, 200, tokenResponse{Token: token})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "something went wrong", err)
		return
	}
	if err = cfg.dbQueries.UpdateRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, 401, "something went wrong", err)
		return
	}
	w.WriteHeader(204)
}

func (cfg *apiConfig) modUserHandler(w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	userJson := &UserReq{}
	err = decoder.Decode(userJson)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userJson.Password)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	user, err := cfg.dbQueries.ModUser(r.Context(), database.ModUserParams{
		ID:             userID,
		Email:          userJson.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, user)
}
