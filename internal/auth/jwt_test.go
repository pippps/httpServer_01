package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pippps/httpServer_01/internal/auth"
)

func signToken(t *testing.T, secret, subject string, expiresIn time.Duration) string {
	t.Helper()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   subject,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("signToken() failed: %v", err)
	}
	return token
}

func TestMakeJWT(t *testing.T) {
	// The token embeds IssuedAt, so its exact string can't be predicted.
	// Instead, check that the token validates back to the same user ID.
	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "valid token",
			userID:      uuid.New(),
			tokenSecret: "secret",
			expiresIn:   time.Hour,
			wantErr:     false,
		},
		{
			name:        "short expiry",
			userID:      uuid.New(),
			tokenSecret: "another-secret",
			expiresIn:   time.Minute,
			wantErr:     false,
		},
		{
			name:        "nil user ID",
			userID:      uuid.Nil,
			tokenSecret: "secret",
			expiresIn:   time.Hour,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got == "" {
				t.Fatal("MakeJWT() returned an empty token")
			}

			id, err := auth.ValidateJWT(got, tt.tokenSecret)
			if err != nil {
				t.Fatalf("ValidateJWT() on MakeJWT() token failed: %v", err)
			}
			if id != tt.userID {
				t.Errorf("ValidateJWT() = %v, want %v", id, tt.userID)
			}
		})
	}

}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	const secret = "secret"

	tests := []struct {
		name        string
		tokenSecret string
		tokenString string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid token",
			tokenSecret: secret,
			tokenString: signToken(t, secret, userID.String(), time.Hour),
			want:        userID,
			wantErr:     false,
		},
		{
			name:        "wrong secret",
			tokenSecret: "wrong-secret",
			tokenString: signToken(t, secret, userID.String(), time.Hour),
			want:        uuid.UUID{},
			wantErr:     true,
		},
		{
			name:        "expired token",
			tokenSecret: secret,
			tokenString: signToken(t, secret, userID.String(), -time.Hour),
			want:        uuid.UUID{},
			wantErr:     true,
		},
		{
			name:        "malformed token",
			tokenSecret: secret,
			tokenString: "not.a.jwt",
			want:        uuid.UUID{},
			wantErr:     true,
		},
		{
			name:        "empty token",
			tokenSecret: secret,
			tokenString: "",
			want:        uuid.UUID{},
			wantErr:     true,
		},
		{
			name:        "subject is not a UUID",
			tokenSecret: secret,
			tokenString: signToken(t, secret, "not-a-uuid", time.Hour),
			want:        uuid.UUID{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}

}
