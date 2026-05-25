package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cli-pokedex/internal/commands"
)

// Creating input store
var inputText []string

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Pokedex!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputText = cleanInput(input)
		val, ok := commands.Commands[inputText[0]]
		if !ok {
			fmt.Println("Your command is:", inputText)
			continue
		}
		if err := val.Callback(inputText); err != nil {
			fmt.Println("Coudn't complete your request. Please try again")
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Unexpected error occured\nClosing the application...")
			commands.Commands["exit"].Callback(inputText)
		}
	}
}


// Helper Function

func cleanInput(text string) []string {
	return strings.Split(strings.ToLower(text), " ")
}
