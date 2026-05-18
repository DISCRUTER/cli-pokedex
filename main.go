package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
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
		if ok {
			val.callback()
		} else {
			fmt.Printf("You command is %s\n", cleanedInput[0])
		}
	}
}
