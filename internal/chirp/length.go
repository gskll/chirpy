package chirp

import "fmt"

func ValidateLength(body string) error {
	if len(body) > 140 {
		return fmt.Errorf("Chirp is too long. Max 140 chars. Actual: %d", len(body))
	}
	return nil
}
