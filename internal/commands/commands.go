package commands

// Command struct
type cliCommand struct {
	Name        string
	Description string
	Usage       string
	Callback    func([]string) error
}

var Commands map[string]cliCommand

func init() {
	Commands = map[string]cliCommand{
		"catch": {
			Name:        "catch",
			Description: "Catch the specified pokemon.",
			Usage:       "catch <pokemon-name>",
			Callback:    catchCommand,
		},
		"explore": {
			Name:        "explore",
			Description: "List of the pokemons encountered in the area.",
			Usage:       "explore <area-name>",
			Callback:    exploreCommand,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect all the caught pokemon.",
			Usage:       "inspect <pokemon-name>",
			Callback:    inspectCommand,
		},
		"map": {
			Name:        "map",
			Description: "List 20 next location areas | Go to next 20",
			Usage:       "map",
			Callback:    mapCommand,
		},
		"mapb": {
			Name:        "mapb",
			Description: "List 20 previous location areas | Go to previous 20",
			Usage:       "mapb",
			Callback:    mapbCommand,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all the pokemon caught.",
			Usage:       "pokedex",
			Callback:    pokedexCommand,
		},
		"help": {
			Name:        "help",
			Description: "Display a help message",
			Usage:       "help",
			Callback:    helpCommand,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Usage:       "exit",
			Callback:    exitCommand,
		},
	}
}
