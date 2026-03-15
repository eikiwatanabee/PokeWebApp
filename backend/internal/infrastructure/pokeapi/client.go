package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
)

const (
	baseURL     = "https://pokeapi.co/api/v2"
	maxPokemon  = 898 // Gen 1-8
	httpTimeout = 10 * time.Second
)

type pokeAPIResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// Client implements domain.service.PokemonFetcher.
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: httpTimeout},
	}
}

func (c *Client) FetchRandom(ctx context.Context) (*valueobject.PokemonInfo, error) {
	id := rand.Intn(maxPokemon) + 1
	return c.FetchByID(ctx, id)
}

func (c *Client) FetchByID(ctx context.Context, id int) (*valueobject.PokemonInfo, error) {
	url := fmt.Sprintf("%s/pokemon/%d", baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokeapi returned status %d", resp.StatusCode)
	}

	var apiResp pokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	types := make([]string, len(apiResp.Types))
	for i, t := range apiResp.Types {
		types[i] = t.Type.Name
	}

	return &valueobject.PokemonInfo{
		PokemonID: apiResp.ID,
		Name:      apiResp.Name,
		SpriteURL: apiResp.Sprites.FrontDefault,
		Types:     types,
	}, nil
}
