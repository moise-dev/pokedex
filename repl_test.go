package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/moise-dev/pokedex/internal/pokecache"
)

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

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(1 * time.Second)
	fmt.Println(cache)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
