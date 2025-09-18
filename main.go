package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/moise-dev/pokedex/internal/api"
	"github.com/moise-dev/pokedex/internal/pokecache"
)

type Config struct {
	prev string
	next string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache, ...string) error
}

func cleanInput(text string) []string {
	return strings.Split(strings.TrimSpace(strings.ToLower(text)), " ")

}
func commandExit(c *Config, cache *pokecache.Cache, args ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, cache *pokecache.Cache, args ...string) error {
	fmt.Println(`Welcome to the Pokedex!
	Usage:

	help:  Displays a help message
	map:   Display 20 locations
	mapb:  Display previous 20 locations
	exit:  Exit the Pokedex`)
	return nil

}

func commandMapGeneric(c *Config, cache *pokecache.Cache, url string) error {
	response, err := api.GetLocation(url, cache)
	if err != nil {
		return err
	}

	c.prev = ""
	c.next = ""
	if response.Prev != nil {
		c.prev = *response.Prev
	}

	if response.Next != nil {
		c.next = *response.Next
	}

	for _, city := range response.Results {
		fmt.Println(city.Name)
	}

	return nil
}

func commandMapBack(c *Config, cache *pokecache.Cache, args ...string) error {
	err := commandMapGeneric(c, cache, c.prev)
	return err

}

func commandMapNext(c *Config, cache *pokecache.Cache, args ...string) error {
	err := commandMapGeneric(c, cache, c.next)
	return err

}

func commandExplore(c *Config, cache *pokecache.Cache, locationName ...string) error {

	if len(locationName) != 2 {
		return errors.New("no location or too many locations provided")
	}

	data, err := api.GetPokemonInLocation(locationName[1], cache)

	for _, encounter := range data.Encounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return err
}

func commandCatch(c *Config, cache *pokecache.Cache, pokemonName ...string) error {
	experience, err := api.CatchPokemon(pokemonName[1], cache)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName[1])
	caughtChance := rand.Intn(experience)
	if caughtChance > experience/2 {
		fmt.Printf("%s was caught!\n", pokemonName[1])

	} else {
		fmt.Printf("%s escaped!\n", pokemonName[1])

	}

	return nil
}

func main() {
	var config Config
	cache := pokecache.NewCache(7 * time.Second)

	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "exit program",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "usage help",
			callback:    commandHelp,
		},

		"map": {
			name:        "map",
			description: "get 20 locations",
			callback:    commandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "map back",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "return pokemons in a certain area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "try to catch a pokemon",
			callback:    commandCatch,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for !scanner.Scan() {
			break
		}
		text := scanner.Text()
		cleanText := cleanInput(text)

		command, ok := commands[cleanText[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(&config, &cache, cleanText...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

}
