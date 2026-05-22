package main

import (
	"bufio"
	"fmt"
	"os"

	"cli-pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

type config struct {
	next     string
	previous string
}

var mapConfig config

var cache pokecache.Cache

func init() {
	commands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "List 20 next location areas | Go to next 20",
			callback:    mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "List 20 previous location areas | Go to previous 20",
			callback:    mapbCommand,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    helpCommand,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exitCommand,
		},
	}

	mapConfig = config{
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		previous: "",
	}

	cache = *pokecache.NewCache()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Pokedex!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		val, ok := commands[cleanedInput[0]]
		if !ok {
			fmt.Printf("You command is %s\n", cleanedInput[0])
			continue
		}
		if err := val.callback(); err != nil {
			fmt.Errorf("%v", err)
		}
	}
}
