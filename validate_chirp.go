package main

import (
	"strings"
)

func profanity(s string) string {
	out := strings.Split(s, " ")
	for i, v := range out {
		if strings.ToLower(v) == "kerfuffle" || strings.ToLower(v) == "sharbert" || strings.ToLower(v) == "fornax" {
			out[i] = "****"
		}
	}

	return strings.Join(out, " ")
}

func isValidChirpBody(s string) bool {
	const maxChirpLength = 140
	if len(s) > maxChirpLength {
		return false
	}
	return true
}
