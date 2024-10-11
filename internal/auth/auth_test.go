package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "mySecurePassword123", false},
		{"Empty password", "", false},
		{"Long password", strings.Repeat("a", 72), false}, // bcrypt has a max input length of 72 bytes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (hash == "" || hash == tt.password) {
				t.Errorf("HashPassword() returned invalid hash for password %q", tt.password)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{"Correct password", password, hash, false},
		{"Incorrect password", "wrongPassword", hash, true},
		{"Empty password", "", hash, true},
		{"Empty hash", password, "", true},
		{"Invalid hash", password, "invalid_hash", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
