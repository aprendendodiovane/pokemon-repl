package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var cliName string = "pokedex-cli"

type cliCommand struct {
	name           string
	description    string
	callback       func(conf *config)
}

type config struct {
	Next 		string	   
	Previous 	string	   
}

func returnCliCommand() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:           "help",
            description:    "Display this help message",
            callback:       displayHelp,
		},
		"version": {
			name:           "version",
            description:    "Display the version of the Pokédex CLI",
            callback:       displayVersion,
        },
		"exit": {
			name:           "exit",
            description:    "Exit the Pokédex CLI",
            callback:       commandExit,
        },
		"map": {
			name:           "map",
            description:    "Display the next map of the Pokémon world",
            callback:       commandMapNext,
        },
		"mapb": {
			name:           "mapb",
            description:    "Display the previous map of the Pokémon world",
            callback:       commandMapPrevious,
        },
	}
}

func displayHelp(conf *config) {
	fmt.Printf("\n%s - Command-line interface for Pokémon information\n\n", cliName)
    fmt.Println("Usage: pokedex-cli [command] [options]")
    fmt.Println("\nCommands:")
    for _, command := range returnCliCommand() {
        fmt.Printf("  %s - %s\n", command.name, command.description)
    }
}

func displayVersion(conf *config) {
    fmt.Printf("%s - Version 1.0.0\n", cliName)
}

func commandExit(conf *config) {
    os.Exit(0)
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

type locations struct {
	Next 		string	   `json:"next"`
	Previous 	string	   `json:"previous"`
	Results 	[]location `json:"results"`
}

type location struct {
    Name string `json:"name"`
}

func commandMapNext(conf *config) {
	url := "https://pokeapi.co/api/v2/location-area/"

	if conf.Next != "" {
		url = conf.Next
	}

	res, err := http.Get(url)
	if err != nil {
        fmt.Printf("Error fetching map data: %v\n", err)
        return
    }

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
        fmt.Printf("Error reading map data: %v\n", err)
        return
    }

	var locations locations
	err = json.Unmarshal(body, &locations)
	if err != nil {
        fmt.Printf("Error parsing map data: %v\n", err)
        return
    }

	conf.Next = locations.Next
	conf.Previous = locations.Previous

	fmt.Println("\nMap of Pokémon world:")
	for _, location := range locations.Results {
        fmt.Printf("  - %s\n", location.Name)
    }
}

func commandMapPrevious(conf *config) {
	url := "https://pokeapi.co/api/v2/location-area/"

	if conf.Previous != "" {
		url = conf.Previous
	}

	res, err := http.Get(url)
	if err != nil {
        fmt.Printf("Error fetching map data: %v\n", err)
        return
    }

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
        fmt.Printf("Error reading map data: %v\n", err)
        return
    }

	var locations locations
	err = json.Unmarshal(body, &locations)
	if err != nil {
        fmt.Printf("Error parsing map data: %v\n", err)
        return
    }

	conf.Next = locations.Next
	conf.Previous = locations.Previous

	fmt.Println("\nMap of Pokémon world:")
	for _, location := range locations.Results {
        fmt.Printf("  - %s\n", location.Name)
    }
}

func getPokemonInfo(name string) (string, error) {
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	res, err := http.Get(fullUrl)
	if err!= nil {
        return "", fmt.Errorf("error fetching data: %v", err)
    }

	defer res.Body.Close()

	if res.StatusCode!= http.StatusOK {
        return "", fmt.Errorf("error fetching data: %s", res.Status)
    }

	pokemonBytes, err := io.ReadAll(res.Body)
	if err != nil {
        return "", fmt.Errorf("error parsing data: %v", err)
    }



	var pokemonBuffer bytes.Buffer
	if err := json.Indent(&pokemonBuffer, pokemonBytes, "", " "); err != nil {
        return "", fmt.Errorf("error parsing data: %v", err)
    }

	return pokemonBuffer.String(), nil
}

func cleanInput(input string) string {
	output := strings.TrimSpace(input)
	output = strings.ToLower(output)
	return output
}

func main() {
	config := &config{}

	commands := returnCliCommand()
	scanner := bufio.NewScanner(os.Stdin)
	printPrompt()
	for scanner.Scan() {
	
		name := cleanInput(scanner.Text())
		if command, ok := commands[name]; ok {
			command.callback(config)
		}

        if name == "" {
            printPrompt()
            continue
        }
	}
}