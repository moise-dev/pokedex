package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/moise-dev/pokedex/api"
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
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil

}

func commandMapGeneric(c *Config, url string) error {
	response, err := api.GetLocation(url)
	if err != nil {
		return err
	}

	c.prev = response.Prev
	c.next = response.Next

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
		command.callback(&config)

	}

}
