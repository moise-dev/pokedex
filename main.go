package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Split(strings.TrimSpace(strings.ToLower(text)), " ")

}

func main() {
	fmt.Print("Hello, World!")

}
