package main

import (
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Split(strings.ToLower(text), " ")
}

func helpCommand() error {
	fmt.Println("Usage\n ")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func exitCommand() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
