package main

import (
	"testing"
)

func TestCensor(t *testing.T) {
	type test struct {
		input  string
		output string
	}
	tests := []test{
		{
			input:  "I had a really rough day today. This kerfuffle is annoying.",
			output: "I had a really rough day today. This **** is annoying.",
		},
		{
			input:  "I hear Mastodon is better than Chirpy. Sharbert I need to migrate",
			output: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			input:  "I really need a Kerfuffle to go to bed sooner, Fornax!",
			output: "I really need a **** to go to bed sooner, Fornax!",
		},
		{
			input:  "hehe, here come's the FOrnaX and his kerfuffle, hopefully i will sharberte this kerfuffle in the ground!",
			output: "hehe, here come's the **** and his kerfuffle, hopefully i will sharberte this **** in the ground!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := censor(tc.input)
			if got != tc.output {
				t.Errorf("censor(%q)\n got: %q\nwant: %q", tc.input, got, tc.output)
			}
		})
	}

}
