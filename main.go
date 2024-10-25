package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aprendendodiovane/pokedex-repl/internal/pokecache"
)

var cliName string = "pokedex-cli"

type cliCommand struct {
	name        string
	description string
	callback    func(conf *config, cache pokecache.Cache, name string)
}

type config struct {
	Next     string
	Previous string
	Pokedex 
}

func returnCliCommand() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display this help message",
			callback:    displayHelp,
		},
		"version": {
			name:        "version",
			description: "Display the version of the Pokédex CLI",
			callback:    displayVersion,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokédex CLI",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next map of the Pokémon world",
			callback:    commandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous map of the Pokémon world",
			callback:    commandMapPrevious,
		},
		"explore": {
			name:        "explore",
            description: "Explore a Pokémon location",
            callback:    commandExplore,
		},
		"catch" : {
			name:        "catch",
            description: "Catch a Pokémon",
            callback:    commandCatch,
		},
		"inspect" : {
			name:        "inspect",
            description: "Inspect a Pokémon",
            callback:    commandInspect,
		},
		"pokedex" : {
			name:        "pokedex",
            description: "List all Pokemon on Pokedex",
            callback:    commandPokedex,
		},
	}
}

func displayHelp(conf *config, cache pokecache.Cache, name string) {
	fmt.Printf("\n%s - Command-line interface for Pokémon information\n\n", cliName)
	fmt.Println("Usage: pokedex-cli [command] [options]")
	fmt.Println("\nCommands:")
	for _, command := range returnCliCommand() {
		fmt.Printf("  %s - %s\n", command.name, command.description)
	}
}

func displayVersion(conf *config, cache pokecache.Cache, name string) {
	fmt.Printf("%s - Version 1.0.0\n", cliName)
}

func commandExit(conf *config, cache pokecache.Cache, name string) {
	os.Exit(0)
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

type locations struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

type location struct {
	Name string `json:"name"`
}

type Explore struct {
	Name                 string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type CatchPokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

type Pokedex struct {
	Pokemons map[string]CatchPokemon
}

func commandPokedex(config *config, cache pokecache.Cache, name string) {
	if len(config.Pokedex.Pokemons) == 0 {
		fmt.Println("No Pokémon found on the Pokedex.")
        return
	}

	fmt.Println("Your Pokedex:")
	for name := range config.Pokedex.Pokemons {
        fmt.Printf(" - %s\n", name)
    }
}

func commandInspect(config *config, cache pokecache.Cache, name string) {
    if len(config.Pokedex.Pokemons) == 0 {
        fmt.Println("You haven't caught any Pokémon yet.")
        return
    }

	pokemonName := strings.ToLower(name)
	pokemon, ok := config.Pokedex.Pokemons[pokemonName]
	if !ok {
        fmt.Printf("No Pokémon found with the name '%s'.\n", name)
        return
    }

    fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
        fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
    }
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
        fmt.Printf(" -%s\n", t.Type.Name)
    }
}

func commandCatch(config *config, cache pokecache.Cache, name string) {
	url := "https://pokeapi.co/api/v2/pokemon/"+ name + "/"

	res, err := http.Get(url)
	if err!= nil {
        fmt.Printf("Error fetching pokemon data: %v\n", err)
        return
    }

	if res.StatusCode > 299 {
        fmt.Printf("Error fetching pokemon data: status code %d\n", res.StatusCode)
        return
    }

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err!= nil {
        fmt.Printf("Error reading data: %v\n", err)
        return
    }

	var pokemonToCatch CatchPokemon
	err = json.Unmarshal(body, &pokemonToCatch)
	if err!= nil {
        fmt.Printf("Error parsing data: %v\n", err)
        return
    }

	randomNumber := rand.IntN(pokemonToCatch.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonToCatch.Name)
	if randomNumber >= pokemonToCatch.BaseExperience/2 {
		fmt.Printf("Congratulations! You caught a %s!\n", pokemonToCatch.Name)
		config.Pokedex.Pokemons[pokemonToCatch.Name] = pokemonToCatch
	} else {
		fmt.Printf("You didn't catch %s.\n", pokemonToCatch.Name)
		fmt.Printf("Good luck next time!\n")
	}
}

func commandExplore(conf *config, cache pokecache.Cache, name string) {
	url := "https://pokeapi.co/api/v2/location-area/" + name + "/"

	var explore Explore

	item, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(item, &explore)
		if err != nil {
			fmt.Printf("Error parsing  data: %v\n", err)
			return
		}
	}

	res, err := http.Get(url)
	if err != nil {
        fmt.Printf("Error fetching location data: %v\n", err)
        return
    }

	if res.StatusCode > 299 {
		fmt.Printf("Error fetching location data: status code %d\n", res.StatusCode)
        return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err!= nil {
        fmt.Printf("Error reading location data: %v\n", err)
        return
    }

	cache.Add(url, body)
	
	err = json.Unmarshal(body, &explore)
	if err!= nil {
        fmt.Printf("Error parsing location data: %v\n", err)
        return
    }

	fmt.Printf("Exploring %s...\n", name)
	fmt.Printf("Found pokemon :\n")
	for _, encounter := range explore.PokemonEncounters {
        fmt.Printf(" - %s\n", encounter.Pokemon.Name)
    }

}

func commandMapNext(conf *config, cache pokecache.Cache, name string) {
	url := "https://pokeapi.co/api/v2/location-area/"

	if conf.Next != "" {
		url = conf.Next
	}

	var locations locations

	item, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(item, &locations)
		if err != nil {
			fmt.Printf("Error parsing map data: %v\n", err)
			return
		}
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

	cache.Add(url, body)

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

func commandMapPrevious(conf *config, cache pokecache.Cache, name string) {
	url := "https://pokeapi.co/api/v2/location-area/"

	if conf.Previous != "" {
		url = conf.Previous
	}

	var locations locations

	item, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(item, &locations)
		if err != nil {
			fmt.Printf("Error parsing map data: %v\n", err)
			return
		}
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

	cache.Add(url, body)

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
	if err != nil {
		return "", fmt.Errorf("error fetching data: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
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
	config.Pokedex.Pokemons = make(map[string]CatchPokemon)
	interval := 10 * time.Second
	cache := pokecache.NewCache(interval)

	commands := returnCliCommand()
	scanner := bufio.NewScanner(os.Stdin)
	printPrompt()
	for scanner.Scan() {
		c := cleanInput(scanner.Text())
		com := strings.TrimSpace(c)
		parts := strings.Split(com, " ")
		if command, ok := commands[parts[0]]; ok {
			if len(parts) == 1 {
				command.callback(config, cache, "")
			} else {
				command.callback(config, cache, parts[1])
			}
		}

		if parts[0] == "" {
			printPrompt()
			continue
		}
	}
}
