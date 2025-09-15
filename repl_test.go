package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world   ",
			expected: []string{"hello", "world"},
		},

		{
			input:    "hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "     hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    " helloworld   ",
			expected: []string{"helloworld"},
		},
		{

			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{

			input:    "",
			expected: []string{""},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("len mismatch. Expected %d, got %d. Expected %s, got %s", len(c.expected), len(actual), c.expected, actual)
		}

		for i, word := range actual {
			if word != c.expected[i] {
				t.Errorf("word mismatch. Expected %s, got %s", word, c.expected[i])
				t.Errorf("word mismatch. Expected %s, got %s. Expected %s, got %s", word, c.expected[i], c.expected, actual)
			}
		}

	}
}
