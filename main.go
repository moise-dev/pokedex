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

type MapMovement struct {
	prev string
	next string
}

type Pokedex map[string]api.PokemonInfo

type App struct {
	mapmove MapMovement
	cache   pokecache.Cache
	pokedex Pokedex
}

type cliCommand struct {
	name        string
	description string
	callback    func(*App, ...string) error
}

func cleanInput(text string) []string {
	return strings.Split(strings.TrimSpace(strings.ToLower(text)), " ")

}
func commandExit(app *App, args ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(app *App, args ...string) error {
	fmt.Println(`Welcome to the Pokedex!
	Usage:

	help:     Displays a help message
	catch:    Try to catch a pokemon
	explore:  See which pokemons are in an area
	map:      Display 20 locations
	mapb:     Display previous 20 locations
	exit:     Exit the Pokedex`)
	return nil

}

func commandMapGeneric(app *App, url string) error {
	response, err := api.GetLocation(url, &app.cache)
	if err != nil {
		return err
	}

	app.mapmove.prev = ""
	app.mapmove.next = ""
	if response.Prev != nil {
		app.mapmove.prev = *response.Prev
	}

	if response.Next != nil {
		app.mapmove.next = *response.Next
	}

	for _, city := range response.Results {
		fmt.Println(city.Name)
	}

	return nil
}

func commandMapBack(app *App, args ...string) error {
	err := commandMapGeneric(app, app.mapmove.prev)
	return err

}

func commandMapNext(app *App, args ...string) error {
	err := commandMapGeneric(app, app.mapmove.next)
	return err

}

func commandExplore(app *App, locationName ...string) error {

	if len(locationName) != 2 {
		return errors.New("no location or too many locations provided")
	}

	data, err := api.GetPokemonInLocation(locationName[1], &app.cache)
	if err != nil {
		if err.Error() == "not found" {
			fmt.Printf("\n")
			return nil
		}
		return err
	}

	for _, encounter := range data.Encounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(app *App, pokemonName ...string) error {
	name := pokemonName[1]
	pokemonInfo, err := api.CatchPokemon(name, &app.cache)
	if err != nil {
		if err.Error() == "not found" {
			fmt.Printf("\n")
			return nil
		}
		return err
	}

	experience := pokemonInfo.BaseExperience

	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	caughtChance := rand.Intn(experience)
	if caughtChance > experience/2 {
		fmt.Printf("%s was caught!\n", name)
		app.pokedex[name] = pokemonInfo

	} else {
		fmt.Printf("%s escaped!\n", name)

	}

	return nil
}

func commandInspect(app *App, pokemonName ...string) error {
	entry, found := app.pokedex[pokemonName[1]]
	if found == false {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemonName[1])
	fmt.Printf("Height: %d\n", entry.Height)
	fmt.Printf("Weight: %d\n", entry.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range entry.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")

	for _, stat := range entry.Types {
		fmt.Printf("  -%s\n", stat.Type.Name)
	}
	return nil
}

func commandPokedex(app *App, args ...string) error {
	if len(app.pokedex) == 0 {
		fmt.Printf("Your Pokedex is empty.\n")
		return nil
	}

	fmt.Printf("Your Pokedex:\n")
	for _, entry := range app.pokedex {
		fmt.Printf("  - %s\n", entry.Name)
	}

	return nil

}

func main() {
	app := App{
		mapmove: MapMovement{},
		cache:   pokecache.NewCache(7 * time.Second),
		pokedex: make(Pokedex),
	}

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
		"inspect": {
			name:        "inspect",
			description: "display a pokedex entry",
			callback:    commandInspect,
		},

		"pokedex": {
			name:        "pokedex",
			description: "display all the caught pokemon's names",
			callback:    commandPokedex,
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

		err := command.callback(&app, cleanText...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

}
