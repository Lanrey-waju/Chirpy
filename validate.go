package main

import (
	"strings"
)

func removeProfaneWords(profaneWords map[string]struct{}, body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if _, ok := profaneWords[lowered]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
