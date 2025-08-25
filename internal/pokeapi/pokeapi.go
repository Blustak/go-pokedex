package pokeapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Blustak/go-pokedex/internal/pokecache"
)

const baseUrl = "https://pokeapi.co/api/v2"
const page_limit = 20

type PokeApiRequest interface {
	GetFullUrl() string
}

type LocationAreaNamesRequest struct {
	Page int
}


func (l LocationAreaNamesRequest) GetFullUrl() string {
	var postfix string
	if offset := l.Page * 20; offset > 0 {
		postfix = fmt.Sprintf("?limit=%d&offset=%d", page_limit, offset)
	}

	return fmt.Sprintf("%s/location-area%s", baseUrl, postfix)
}

type LocationAreaExploreRequest struct {
    Name string
}

func (l LocationAreaExploreRequest) GetFullUrl() string {
    return fmt.Sprintf("%s/location-area/%s", baseUrl, l.Name)
}

type PokemonRequest struct {
    Name string
}

func (p PokemonRequest) GetFullUrl() string {
    return fmt.Sprintf("%s/pokemon/%s", baseUrl, p.Name)
}

func Get(req PokeApiRequest, cache *pokecache.Pokecache) ([]byte, error) {
	fullUrl := req.GetFullUrl()
	if cacheRes, ok := cache.Get(fullUrl); ok {
        if cacheRes == nil {
            return nil, fmt.Errorf("error resource does not exist")
        }
		return cacheRes, nil
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching response: %w", err)
	}
    if res.StatusCode == 404 {
        cache.Add(fullUrl, nil) 
    }
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("error bad status code %s", res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	cache.Add(fullUrl, body)
	return body, nil
}
