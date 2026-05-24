package main

import (
	"bufio"
	"fmt"
	"os"

	"cli-pokedex/internal/pokecache"
)

// Command struct
type cliCommand struct {
	name        string
	description string
	usage       string
	callback    func() error
}

var commands map[string]cliCommand

// Creating input store
var inputText []string

// Declaring config for url nagivation
type config struct {
	next     string
	previous string
}

var collection map[string]pokemon

var mapConfig config

// Creating cache type
var cache pokecache.Cache

func init() {
	commands = map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Catch the specified pokemon.",
			usage:       "catch <pokemon-name>",
			callback:    catchCommand,
		},
		"explore": {
			name:        "explore",
			description: "List of the pokemons encountered in the area.",
			usage:       "explore <area-name>",
			callback:    exploreCommand,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect all the caught pokemon.",
			usage:       "inspect <pokemon-name>",
			callback:    inspectCommand,
		},
		"map": {
			name:        "map",
			description: "List 20 next location areas | Go to next 20",
			usage:       "map",
			callback:    mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "List 20 previous location areas | Go to previous 20",
			usage:       "mapb",
			callback:    mapbCommand,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the pokemon caught.",
			usage:       "pokedex",
			callback:    pokedexCommand,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			usage:       "help",
			callback:    helpCommand,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			usage:       "exit",
			callback:    exitCommand,
		},
	}

	// Initializing map config
	mapConfig = config{
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		previous: "",
	}

	// Configuring cache
	pokecache.SetCacheDuration(10)
	cache = *pokecache.NewCache()

	// Pokemon Collection
	collection = make(map[string]pokemon)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Pokedex!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputText = cleanInput(input)
		val, ok := commands[inputText[0]]
		if !ok {
			fmt.Println("Your command is:", inputText)
			continue
		}
		if err := val.callback(); err != nil {
			fmt.Errorf("%v", err)
		}
	}
}
