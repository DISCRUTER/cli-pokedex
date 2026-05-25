package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	// "cli-pokedex/src/main"
	"cli-pokedex/internal/pokecache"
)

// Declaring config for url nagivation
type config struct {
	next     string
	previous string
}

var mapConfig config

// Pokedex
type pokedex struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

// Pokemon Location
type pokemonLocation struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Pokemon
type pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

var collection map[string]pokemon

// Creating cache type
var cache pokecache.Cache

func init() {
	// Initializing map config
	mapConfig = config{
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		previous: "",
	}
	// Initializing map config
	mapConfig = config{
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		previous: "",
	}
	// Configuring cache
	pokecache.SetCacheDuration(10)
	cache = *pokecache.NewCache()
}

func catchCommand(inputText []string) error {
	if len(inputText) < 2 {
		return fmt.Errorf("Location name not provided.\nUsage: %v", Commands["catch"].Usage)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", inputText[1])

	// Fetch pokemon data
	URL := "https://pokeapi.co/api/v2/pokemon/" + inputText[1]
	data, err := getApiData(URL)
	if err != nil {
		return err
	}

	// Unmarshal the data
	var pokemonData pokemon
	if err = json.Unmarshal(data, &pokemonData); err != nil {
		return nil
	}

	// Genrating odd
	r := rand.New(rand.NewSource(time.Now().UnixNano() * int64(pokemonData.BaseExperience)))
	if r.Float32() > 0.5 {
		fmt.Printf("%s was caught!\n", pokemonData.Name)
		collection[pokemonData.Name] = pokemonData
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
	}

	return nil
}

func exploreCommand(inputText []string) error {
	// Check if argument is given
	if len(inputText) < 2 {
		return fmt.Errorf("Location name not provided.\nUsage: %v", Commands["explore"].Usage)
	}

	// Fetch pokemon names
	URL := "https://pokeapi.co/api/v2/location-area/" + inputText[1]

	// Fetch map data
	data, err := getApiData(URL)
	if err != nil {
		return err
	}
	// Unmarshal the data
	var pokemonInfo pokemonLocation
	if err = json.Unmarshal(data, &pokemonInfo); err != nil {
		return err
	}
	// Print the result
	fmt.Println("Potential Pokemon encounters...")
	for _, encounter := range pokemonInfo.PokemonEncounters {
		println(encounter.Pokemon.Name)
	}

	return nil
}

func inspectCommand(inputText []string) error {
	// Check if pokemon name exist
	if len(inputText) < 2 {
		return fmt.Errorf("Location name not provided.\nUsage: %v", Commands["inspect"].Usage)
	}

	// Fetch pokemon data
	data, exist := collection[inputText[1]]
	if !exist {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// Printing data
	fmt.Printf("Name: %s\nBase Experience: %d\nHeight: %d\nWeight: %d\n", data.Name, data.BaseExperience, data.Height, data.Weight)
	fmt.Println("Types:")
	for _, val := range data.Types {
		fmt.Println("  -", val.Type.Name)
	}

	return nil
}

func mapCommand(inputText []string) error {
	// Check if next exist
	if mapConfig.next == "" {
		fmt.Println("That's the end.\nUse `mapb` to go back.")
		return fmt.Errorf("Next link not found.")
	}

	// Fetch map data
	data, err := getApiData(mapConfig.next)
	if err != nil {
		return err
	}
	// Unmarshal map data
	if err = printMapData(data); err != nil {
		return err
	}

	return nil
}

func mapbCommand(inputText []string) error {
	// Check if previous exist
	if mapConfig.previous == "" {
		fmt.Println("This is the start.\nUse `map` to go forward.")
		return fmt.Errorf("Previous link not found.")
	}

	// Fetch map data
	data, err := getApiData(mapConfig.previous)
	if err != nil {
		return err
	}
	// Unmarshal map data
	if err = printMapData(data); err != nil {
		return err
	}

	return nil
}

func pokedexCommand(inputText []string) error {
	if len(collection) < 1 {
		fmt.Println("No pokemon found. Try catching one.")
		return fmt.Errorf("No pokemon found.")
	}
	for name := range collection {
		println("  -", name)
	}
	return nil
}

func helpCommand(inputText []string) error {
	fmt.Println("Usage\n ")
	for _, command := range Commands {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func exitCommand(inputText []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	cache.Stop()
	os.Exit(0)
	return nil
}

// Helper functions

func getApiData(URL string) ([]byte, error) {
	// Checking cache
	data, exist := cache.Get(URL)
	if exist {
		return data, nil
	}
	// Sending Http request
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Checking Status Codes
	if res.StatusCode != http.StatusOK {
		fmt.Println("Something went wrong!")
		fmt.Println("Tips: Ensure you have typed a valid pokemon name.")
		return nil, fmt.Errorf("Didn't recieve 200 status code.")
	}
	// Converting data into bytes
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// Updating Cache
	cache.Add(URL, data)
	return data, nil
}

func printMapData(data []byte) error {
	// Unmarshal the bytes
	var pokedex_map pokedex
	if err := json.Unmarshal(data, &pokedex_map); err != nil {
		return err
	}

	// Printing reuslts
	for _, val := range pokedex_map.Results {
		fmt.Printf("%s\n", val.Name)
	}

	// Configuring next and prev urls
	mapConfig.next = pokedex_map.Next
	mapConfig.previous = pokedex_map.Previous

	return nil
}
