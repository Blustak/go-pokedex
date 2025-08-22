package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name     string
	desc     string
	callback func() error
}


var registry map[string]cliCommand

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
	var buf []string
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
			if err := cmd.callback(); err != nil {
				fmt.Printf("Error: %v", err)

			}
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	var err error
	return err
}

func commandHelp() error {
    fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
    for _, v := range registry {
        fmt.Printf("%s: %s\n",v.name,v.desc)
    }
    var err error
    return err
}
