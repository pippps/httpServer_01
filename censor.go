package main

import "strings"

func censor(msg string) string {
	censoredWord := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")
	var censoredMsg string
	for i, word := range words {
		for _, cens := range censoredWord {
			if strings.ToLower(word) == cens {
				words[i] = "****"
			}
		}
		censoredMsg += words[i]
		if i < len(words)-1 {
			censoredMsg += " "
		}
	}

	return censoredMsg
}
