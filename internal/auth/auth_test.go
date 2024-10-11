package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	// Test MakeJWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("MakeJWT returned an error: %v", err)
	}
	if token == "" {
		t.Error("MakeJWT returned an empty token")
	}

	// Test ValidateJWT with valid token
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("ValidateJWT returned an error for a valid token: %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("ValidateJWT returned incorrect userID. Got %v, want %v", parsedUserID, userID)
	}

	// Test ValidateJWT with invalid secret
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("ValidateJWT did not return an error for an invalid secret")
	}

	// Test ValidateJWT with expired token
	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		Subject:   userID.String(),
	}).SignedString([]byte(tokenSecret))
	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT did not return an error for an expired token")
	}

	// Test ValidateJWT with invalid issuer
	invalidIssuerToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:  "not-chirpy",
		Subject: userID.String(),
	}).SignedString([]byte(tokenSecret))
	_, err = ValidateJWT(invalidIssuerToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT did not return an error for an invalid issuer")
	}

	// Test ValidateJWT with invalid user ID
	invalidUserIDToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:  "chirpy",
		Subject: "not-a-uuid",
	}).SignedString([]byte(tokenSecret))
	_, err = ValidateJWT(invalidUserIDToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT did not return an error for an invalid user ID")
	}

	// Test ValidateJWT with invalid signing method
	invalidMethodToken, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Issuer:  "chirpy",
		Subject: userID.String(),
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)
	_, err = ValidateJWT(invalidMethodToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT did not return an error for an invalid signing method")
	}
}
