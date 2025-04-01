package main

import (
	"testing"
)

func TestReplaceBadWords(t *testing.T) {
	cases := []struct {
		input 	 string
		expected string
	}{
		{
			input: "This is a kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			input: "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			input: "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
		{
			input: "I really need a kerfuffle to go to bed sooner, Fornax!",
			expected: "I really need a **** to go to bed sooner, Fornax!",
		},
	}

	for _, c := range cases {
		actual := replaceBadWords(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual, c.expected)
			continue
		}
		if actual != c.expected {
			t.Errorf("Content dont' match")
		}
	}
}
