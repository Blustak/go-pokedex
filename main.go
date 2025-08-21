package main

import (
    "fmt"
    "os"
    "strings"
    "bufio"
)

func main() {
    var buf []string
    reader := bufio.NewReader(os.Stdin)
    scanner := bufio.NewScanner(reader)
    for {
        buf = []string{}
        fmt.Print("Pokedex > ")
        scanner.Scan()
        buf = cleanInput(scanner.Text())
        if len(buf) > 0 {
            fmt.Printf("Your command was: %s\n", buf[0])
        } else {
            fmt.Println("No command given")
        }
    }
}

func cleanInput(text string) []string {
    return strings.Fields(strings.ToLower(text))
}
