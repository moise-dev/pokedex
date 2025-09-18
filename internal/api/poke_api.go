package api

import (
	"encoding/json"
	"errors"
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

type PokemonInfo struct {
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
	BaseExperience int           `json:"base_experience"`
	Name           string        `json:"name"`
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
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
			fmt.Printf("Location %s not found", placeName)
			return PokemonLocation{}, errors.New("not found")
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

func CatchPokemon(pokemonName string, cache *pokecache.Cache) (PokemonInfo, error) {
	fullURL := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	var respBytes []byte

	respBytes, cacheHit := cache.Get(fullURL)

	if cacheHit == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("http connection error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Pokemon %s not found", pokemonName)
			return PokemonInfo{}, errors.New("not found")
		}

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("cannot read http data: %w", err)
		}
		cache.Add(fullURL, respBytes)
	}

	var response PokemonInfo
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		return PokemonInfo{}, fmt.Errorf("cannot unmarshal data: %w", err)
	}

	return response, nil
}
