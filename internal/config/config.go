package config

import "github.com/Blustak/go-pokedex/internal/pokeapi"

type Config struct {
    Next pokeapi.PokeApiRequest
    Previous pokeapi.PokeApiRequest
    LocationAreaPage int
}
