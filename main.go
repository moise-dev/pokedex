package main

import (
	"bufio"
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Split(strings.TrimSpace(strings.ToLower(text)), " ")

}

func main() {

	scanner := bufio.NewScanner(strings.NewReader(""))

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			text := scanner.Text()
		}
		cleanText := cleanInput(text)
		fmt.Println("Your command was: %s", cleanText[0])
	}

}
