package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Blustak/go-pokedex/internal/config"
	"github.com/Blustak/go-pokedex/internal/pokeapi"
	"github.com/Blustak/go-pokedex/internal/pokecache"
	"github.com/Blustak/go-pokedex/internal/pokemon"
)

const cacheClearInterval time.Duration = 5 * time.Minute

type cliCommand struct {
	name     string
	desc     string
	callback func(c *config.Config) error
}

var registry map[string]cliCommand
var cache pokecache.Pokecache
var caughtPokemon map[string]pokemon.Pokemon
var conf config.Config
var args []string

func main() {
	registry = make(map[string]cliCommand)
    caughtPokemon = make(map[string]pokemon.Pokemon)
	cache = pokecache.NewCache(cacheClearInterval)

	registry["exit"] = cliCommand{
		name:     "exit",
		desc:     "Exit the Pokedex",
		callback: commandExit,
	}

	registry["help"] = cliCommand{
		name:     "help",
		desc:     "Displays a help message",
		callback: commandHelp,
	}

	registry["map"] = cliCommand{
		name:     "map",
		desc:     "Get the next 20 location areas",
		callback: commandMap,
	}

	registry["mapb"] = cliCommand{
		name:     "map",
		desc:     "Get the previous 20 location areas",
		callback: commandMapb,
	}

	registry["explore"] = cliCommand{
		name:     "explore",
		desc:     "explore <area>: show pokemon that can be encountered in <area>",
		callback: commandExplore,
	}

	registry["catch"] = cliCommand{
		name:     "catch",
		desc:     "catch <name>: attempt to catch a pokemon with <name>",
		callback: commandCatch,
	}
    
    registry["inspect"] = cliCommand{
        name:   "inspect",
        desc:   "inspect <name>, where <name> is a pokemon you have caught",
        callback: commandInspect,
    }

    registry["pokedex"] = cliCommand{
        name: "pokedex",
        desc: "list the pokemon you have caught",
        callback: commandPokedex,
    }

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	for {
		args = nil
		fmt.Print("Pokedex > ")
		scanner.Scan()
		args = cleanInput(scanner.Text())
		if len(args) <= 0 {
			fmt.Println("No command given")
		}
		cmd, ok := registry[args[0]]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			args = args[1:]
			if err := cmd.callback(&conf); err != nil {
				fmt.Printf("Error: %s\n", fmt.Errorf("%w", err))

			}
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config.Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range registry {
		fmt.Printf("%s: %s\n", v.name, v.desc)
	}
	return nil
}

func commandMap(c *config.Config) error {
	if err := getMap(c); err != nil {
		return err
	}
	c.LocationAreaPage += 1
	return nil
}

func commandMapb(c *config.Config) error {
	if c.LocationAreaPage <= 0 {
		fmt.Println("you're on the first page")
		return nil
	}
	c.LocationAreaPage -= 1
	if err := getMap(c); err != nil {
		return err
	}
	return nil
}

func getMap(c *config.Config) error {
	req := pokeapi.LocationAreaNamesRequest{
		Page: c.LocationAreaPage,
	}
	buf, err := pokeapi.Get(req, &cache)
	if err != nil {
		return err
	}

	var res struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}
	for _, c := range res.Results {
		fmt.Printf("%s\n", c.Name)
	}
	return nil
}

func commandExplore(c *config.Config) error {
	if len(args) != 1 {
		return fmt.Errorf("error expected 1 arg, got %d", len(args))
	}
	req := pokeapi.LocationAreaExploreRequest{
		Name: args[0],
	}
	buf, err := pokeapi.Get(req, &cache)
	if err != nil {
		return err
	}
	var res struct {
		Encounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	for _, v := range res.Encounters {
		fmt.Printf("%s\n", v.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config.Config) error {
	if len(args) != 1 {
		return fmt.Errorf("error, expected 1 argument, got %d", len(args))
	}
	req := pokeapi.PokemonRequest{
		Name: args[0],
	}
	var poke pokemon.Pokemon
	buf, err := pokeapi.Get(req, &cache)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, &poke); err != nil {
		return err
	}
	if _, ok := caughtPokemon[poke.Name]; ok {
		fmt.Printf("you have already caught %s!\n", poke.Name)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", poke.Name)

	if poke.TryCatch() {
		fmt.Printf("%s was caught!\n", poke.Name)
		caughtPokemon[poke.Name] = poke
	} else {
		fmt.Printf("%s escaped!\n", poke.Name)
	}
    return nil
}

func commandInspect(c *config.Config) error {
	if len(args) != 1 {
		return fmt.Errorf("error, expected 1 argument, got %d", len(args))
	}
    poke, ok := caughtPokemon[args[0]]
    if !ok {
        fmt.Println("you have not caught that pokemon")
        return nil
    }
    poke.Print()
    return nil
}

func commandPokedex(c *config.Config) error {
    if len(caughtPokemon) == 0 {
        fmt.Println("You have not caught any pokemon!")
        return nil
    }
    fmt.Println("Your Pokedex:")
    for k, _ := range caughtPokemon {
        fmt.Printf("  - %s\n", k)
    }
    return nil
}
