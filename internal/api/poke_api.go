package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/moise-dev/pokedex/internal/pokecache"
)

type AvailableLocations struct {
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

type PokemonStats struct {
	BaseExperince int `json:"base_experience"`
}

func GetLocation(fullURL string, cache *pokecache.Cache) (AvailableLocations, error) {
	if fullURL == "" {
		fullURL = "https://pokeapi.co/api/v2/location-area/"
	}

	var respBytes []byte

	respBytes, cacheHit := cache.Get(fullURL)

	if cacheHit == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return AvailableLocations{}, fmt.Errorf("http connection error: %w", err)
		}

		defer resp.Body.Close()

		respBytes, err = io.ReadAll(resp.Body)

		if err != nil {
			return AvailableLocations{}, fmt.Errorf("cannot read http data: %w", err)
		}
		cache.Add(fullURL, respBytes)
	}

	var response AvailableLocations
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		return AvailableLocations{}, fmt.Errorf("cannot unmarshal data: %w", err)
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
			return PokemonLocation{}, fmt.Errorf("http connection error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return PokemonLocation{}, fmt.Errorf("%s not found", placeName)
		}

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return PokemonLocation{}, fmt.Errorf("cannot read http data: %w", err)
		}

		cache.Add(fullURL, respBytes)
	}

	var response PokemonLocation
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		return PokemonLocation{}, fmt.Errorf("cannot unmarshal data: %w", err)
	}

	return response, nil

}

func CatchPokemon(pokemonName string, cache *pokecache.Cache) (int, error) {
	fullURL := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	var respBytes []byte

	respBytes, cacheHit := cache.Get(fullURL)

	if cacheHit == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return -1, fmt.Errorf("http connection error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return -1, fmt.Errorf("%s not found", pokemonName)
		}

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return -1, fmt.Errorf("cannot read http data: %w", err)
		}
		cache.Add(fullURL, respBytes)
	}

	var response PokemonStats
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		return -1, fmt.Errorf("cannot unmarshal data: %w", err)
	}

	return response.BaseExperince, nil
}
