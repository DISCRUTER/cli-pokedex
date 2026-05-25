package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "Hello, World!",
			expected: []string{"hello,", "world!"},
		},
		{
			input:    "Go Pikachu",
			expected: []string{"go", "pikachu"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		output := cleanInput(c.input)
		if len(output) != len(c.expected) {
			t.Fatalf("expected: %v, got: %v", c.expected, output)
		}
		for i := range output {
			word := output[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Fatalf("expected: %v, got: %v", expectedWord, word)
			}
		}
	}
}