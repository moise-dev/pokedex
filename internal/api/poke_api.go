package api

import (
	"encoding/json"
	"errors"
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

type PokemonLocation struct {
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocation(fullURL string, cache *pokecache.Cache) (Map, error) {
	if fullURL == "" {
		fullURL = "https://pokeapi.co/api/v2/location-area/"
	}

	var respBytes []byte

	respBytes, cacheHit := cache.Get(fullURL)

	if cacheHit == false {
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
		return Map{}, errors.New("cannot unmarshal data")
	}

	return response, nil

}

func GetPokemonInLocation(placeName string, cache *pokecache.Cache) (PokemonLocation, error) {
	fullURL := "https://pokeapi.co/api/v2/location-area/" + placeName + "/"

	var respBytes []byte

	respBytes, cacheHit := cache.Get(fullURL)

	if cacheHit == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return PokemonLocation{}, errors.New("http connection error")
		}
		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return PokemonLocation{}, errors.New("cannot read http data")
		}

		cache.Add(fullURL, respBytes)
	}

	var response PokemonLocation
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		return PokemonLocation{}, errors.New("cannot unmarshal data")
	}

	return response, nil

}
