package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type pokedex struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func cleanInput(text string) []string {
	return strings.Split(strings.ToLower(text), " ")
}

func mapCommand() error {
	// Check if next exist
	if mapConfig.next == "" {
		fmt.Println("That's the end.\nUse `mapb` to go back.")
		return fmt.Errorf("Next link not found.")
	}

	// Fetch map data
	data, err := getMapData(mapConfig.next)
	if err != nil {
		return err
	}
	// Unmarshal map data
	if err = printMapData(data); err != nil {
		return err
	}

	return nil
}

func mapbCommand() error {
	// Check if previous exist
	if mapConfig.previous == "" {
		fmt.Println("This is the start.\nUse `map` to go forward.")
		return fmt.Errorf("Previous link not found.")
	}

	// Fetch map data
	data, err := getMapData(mapConfig.previous)
	if err != nil {
		return err
	}
	// Unmarshal map data
	if err = printMapData(data); err != nil {
		return err
	}

	return nil
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
	cache.Stop()
	os.Exit(0)
	return nil
}

// Helper functions

func getMapData(URL string) ([]byte, error) {
	// Checking cache
	data, exist := cache.Get(URL)
	if exist {
		return  data, nil
	}
	// Sending Http request
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
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
