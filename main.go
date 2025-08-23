package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
    "github.com/Blustak/go-pokedex/internal/config"
    "github.com/Blustak/go-pokedex/internal/poke_api"
)

type cliCommand struct {
	name     string
	desc     string
	callback func(*config.Config) error
}


var registry map[string]cliCommand
var api pokeApi.PokeApi

func main() {
    registry = make(map[string]cliCommand)
    registry["exit"] = cliCommand{
        name: "exit",
        desc: "Exit the Pokedex",
        callback: commandExit,
    }
    registry["help"] = cliCommand{
        name:"help",
        desc: "Displays a help message",
        callback: commandHelp,
    }
    registry["map"] = cliCommand{
        name:"map",
        desc:"Gets 20 location areas",
        callback: commandMap,
    }
    registry["mapb"] = cliCommand{
        name:"mapb",
        desc:"Gets the previous 20 location areas",
        callback:commandMapb,
    }
	var buf []string
    var conf config.Config
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	for {
		buf = []string{}
		fmt.Print("Pokedex > ")
		scanner.Scan()
		buf = cleanInput(scanner.Text())
		if len(buf) <= 0 {
			fmt.Println("No command given")
		}
		comm := buf[0]
		cmd, ok := registry[comm]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			if err := cmd.callback(&conf); err != nil {
				fmt.Printf("Error: %v", err)

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
	var err error
	return err
}

func commandHelp(c *config.Config) error {
    fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
    for _, v := range registry {
        fmt.Printf("%s: %s\n",v.name,v.desc)
    }
    var err error
    return err
}

func commandMap(c *config.Config) error {
    if err := getMap(c); err != nil {
        return err
    }
    c.Next()
    return nil
}

func commandMapb(c *config.Config) error {
    if c.Page == 0 {
        fmt.Println("You are on the first page.")
        return nil
    }
    c.Previous()
    if err := getMap(c); err != nil {
        return err
    }
    return nil
}

func getMap(c *config.Config) error {
    c.Endpoint = pokeApi.LocationArea
    if err := api.Get(c); err != nil {
        return err
    }
    var tmp struct {
        Results []struct{Name string `json:"name"`} `json:"Results"`
    }
    if err := api.Unmarshal(&tmp); err != nil {
        return err
    }
    for _, v := range tmp.Results {
        fmt.Printf("%s\n", v.Name)
    }
    return nil
}
