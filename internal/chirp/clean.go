package chirp

import "strings"

var profanities = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func Clean(body string) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		lower := strings.ToLower(word)
		if profanities[lower] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
