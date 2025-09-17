package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/moise-dev/pokedex/internal/pokecache"
)

type Map struct {
	Next    *string `json:"next"`
	Prev    *string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
	}
}

func GetLocation(fullURL string, cache *pokecache.Cache) (Map, error) {
	if fullURL == "" {
		fullURL = "https://pokeapi.co/api/v2/location-area/"
	}

	var respBytes []byte

	respBytes, present := cache.Get(fullURL)

	if present == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return Map{}, errors.New("http connection error")
		}

		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)

		if err != nil {
			return Map{}, errors.New("cannot read http data")
		}
		cache.Add(fullURL, respBytes)
	}

	var response Map
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		fmt.Println(err)
		return Map{}, errors.New("cannot unmarshal data")
	}

	return response, nil

}
