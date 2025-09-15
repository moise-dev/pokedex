package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Map struct {
	Results []struct {
		Name string `json:"name"`
	}
}

func GetLocation() (Map, error) {
	fullURL := "https://pokeapi.co/api/v2/location-area/"
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
