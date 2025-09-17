package main

import (
	"bufio"
	"fmt"
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
	callback    func(*Config) error
}

func cleanInput(text string) []string {
	return strings.Split(strings.TrimSpace(strings.ToLower(text)), " ")

}
func commandExit(c *Config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Println(`Welcome to the Pokedex!
	Usage:

	help:  Displays a help message
	map:   Display 20 locations
	mapb:  Display previous 20 locations
	exit:  Exit the Pokedex`)
	return nil

}

func commandMapGeneric(c *Config, url string) error {
	cache := pokecache.NewCache(7 * time.Second)
	response, err := api.GetLocation(url, &cache)
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

func commandMapBack(c *Config) error {
	err := commandMapGeneric(c, c.prev)
	return err

}

func commandMapNext(c *Config) error {
	err := commandMapGeneric(c, c.next)
	return err

}

func main() {
	var config Config
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
		err := command.callback(&config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

}
