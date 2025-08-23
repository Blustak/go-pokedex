package pokeApi

import (
    "encoding/json"
	"fmt"
	"io"
	"net/http"
    "github.com/Blustak/go-pokedex/internal/config"
)

// To be succesful, I need to:
// Check if I have the response cached already - create a map using url as key?
// Load the response into memory, either through cache or GET
// process the JSON
// return the value, and increment/decrement the page, depending on the
// command

const baseUrl = "https://pokeapi.co/api/v2"
const pageLength = 20

type Endpoint int

const (
    LocationArea = iota
)

func GetEndpoint(e Endpoint) (string, error) {
    switch e {
    case LocationArea:
        return "/location-area", nil
    default:
        return "", fmt.Errorf("Error, %d is not a valid endpoint", int(e))
    }

}

type PokeApi struct {
    Response []byte
}



func (p *PokeApi) Get(conf *config.Config) (error) {
    endPoint, err := GetEndpoint(Endpoint(conf.Endpoint))
    if err != nil {
        return err
    }
    fullURL := baseUrl + endPoint
    if conf.Page > 0 {
        fullURL = fullURL + fmt.Sprintf("?offset=%d", conf.Page*pageLength)
    }
    res, err := http.Get(fullURL)
    if err != nil {
        return fmt.Errorf("error: getting from pokeapi: %w", err)
    }
    defer res.Body.Close()
    if res.StatusCode > 299 {
        return fmt.Errorf("error: response failed with status code: %d and \nbody: %s", res.StatusCode, res.Body)
    }
    p.Response, err = io.ReadAll(res.Body)
    if err != nil {
        p.Response = nil
        return fmt.Errorf("error: reading from response body: %w", err)
    }
    return nil
}

func (p *PokeApi) Unmarshal(target any) error {
    if p == nil {
        return fmt.Errorf("error, pokeapi is nil")
    }
    if p.Response == nil {
        return fmt.Errorf("error: tried to Unmarshal a nil response")
    }
    if err := json.Unmarshal(p.Response, &target); err != nil {
        return fmt.Errorf("error: failed to unmarshal response %s: %w", string(p.Response), err)
    }
    return nil
}
