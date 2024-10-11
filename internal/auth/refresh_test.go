package auth

import (
	"encoding/hex"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	t.Run("Generate valid token", func(t *testing.T) {
		token, err := MakeRefreshToken()
		if err != nil {
			t.Fatalf("MakeRefreshToken() returned an error: %v", err)
		}
		if len(token) != 64 {
			t.Errorf("Expected token length of 64, got %d", len(token))
		}
		if !isHex(token) {
			t.Errorf("Token is not a valid hexadecimal string: %s", token)
		}
	})

	t.Run("Generate multiple unique tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			token, err := MakeRefreshToken()
			if err != nil {
				t.Fatalf("MakeRefreshToken() returned an error: %v", err)
			}
			if tokens[token] {
				t.Errorf("Generated duplicate token: %s", token)
			}
			tokens[token] = true
		}
	})
}

func isHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}
