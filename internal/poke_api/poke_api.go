package internal

import (
	"fmt"
	"io"
	"net/http"
)

// To be succesful, I need to:
// Check if I have the response cached already - create a map using url as key?
// Load the response into memory, either through cache or GET
// process the JSON
// return the value, and increment/decrement the page, depending on the
// command

const baseUrl = "https://pokeapi.co/api/v2"

type Endpoint string

const (
    LocationArea = "/location-area"
)

type PokeApi struct {
    Response []byte
}

func (p *PokeApi) Get(endpoint Endpoint) (error) {
    res, err := http.Get(baseUrl + string(endpoint))
    if err != nil {
        return fmt.Errorf("error getting from pokeapi: %w", err)
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if res.StatusCode > 299 {
        return fmt.Errorf("error: response failed with status code: %d and \nbody: %s", res.StatusCode, body)
    }
    p.Response = body
    return nil
}
