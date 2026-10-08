package auth_test

import (
	"strings"
	"testing"

	"github.com/pippps/httpServer_01/internal/auth"
)

func TestHashPassword(t *testing.T) {
	passwords := []string{"Hello, World!", "HJFEOIU534hfwef", ""}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			hash, err := auth.HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword(%q) returned error: %v", password, err)
			}
			if !strings.HasPrefix(hash, "$argon2id$") {
				t.Errorf("HashPassword(%q) = %q, want argon2id-encoded hash", password, hash)
			}

			other, err := auth.HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword(%q) second call returned error: %v", password, err)
			}
			if hash == other {
				t.Errorf("HashPassword(%q) returned the same hash twice; salt should be random", password)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	const password = "correctPassword123!"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned error: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
		wantErr  bool
	}{
		{"matching password", password, hash, true, false},
		{"wrong password", "wrongPassword", hash, false, false},
		{"empty password", "", hash, false, false},
		{"malformed hash", password, "not-a-hash", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
