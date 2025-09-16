package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Map struct {
	Next    string `json:"next"`
	Prev    string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
	}
}

func GetLocation(fullURL string) (Map, error) {
	if fullURL == "" {
		fullURL = "https://pokeapi.co/api/v2/location-area/"
	}
	resp, err := http.Get(fullURL)
	if err != nil {
		os.Exit(1)
		return Map{}, errors.New("http connection error")
	}
	var response Map
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)

	return response, nil

}
