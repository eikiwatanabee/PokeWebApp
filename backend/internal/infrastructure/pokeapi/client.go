package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

type speciesResponse struct {
	IsLegendary        bool `json:"is_legendary"`
	IsMythical         bool `json:"is_mythical"`
	EvolvesFromSpecies *struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"evolves_from_species"`
}

// cachedSpecies stores species rarity info to avoid repeated API calls.
type cachedSpecies struct {
	rarity valueobject.PokemonRarity
}

// Client implements domain.service.PokemonFetcher.
type Client struct {
	httpClient   *http.Client
	speciesCache sync.Map // map[int]*cachedSpecies
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
	// Fetch basic pokemon data
	pokemonURL := fmt.Sprintf("%s/pokemon/%d", baseURL, id)
	var apiResp pokeAPIResponse
	if err := c.fetchURL(ctx, pokemonURL, &apiResp); err != nil {
		return nil, fmt.Errorf("fetch pokemon: %w", err)
	}

	types := make([]string, len(apiResp.Types))
	for i, t := range apiResp.Types {
		types[i] = t.Type.Name
	}

	// Determine rarity from species data
	rarity := c.determineRarity(ctx, id)

	return &valueobject.PokemonInfo{
		PokemonID: apiResp.ID,
		Name:      apiResp.Name,
		SpriteURL: apiResp.Sprites.FrontDefault,
		Types:     types,
		Rarity:    rarity,
	}, nil
}

func (c *Client) determineRarity(ctx context.Context, id int) valueobject.PokemonRarity {
	// Check cache first
	if cached, ok := c.speciesCache.Load(id); ok {
		return cached.(*cachedSpecies).rarity
	}

	rarity := c.fetchAndDetermineRarity(ctx, id)

	// Cache result
	c.speciesCache.Store(id, &cachedSpecies{rarity: rarity})
	return rarity
}

func (c *Client) fetchAndDetermineRarity(ctx context.Context, id int) valueobject.PokemonRarity {
	speciesURL := fmt.Sprintf("%s/pokemon-species/%d", baseURL, id)
	var species speciesResponse
	if err := c.fetchURL(ctx, speciesURL, &species); err != nil {
		return valueobject.RarityCommon
	}

	// Rule 1: Mythical pokemon → Legendary rarity
	if species.IsMythical {
		return valueobject.RarityLegendary
	}

	// Rule 2: Legendary pokemon → Epic rarity
	if species.IsLegendary {
		return valueobject.RarityEpic
	}

	// Rule 3: Check evolution stage
	if species.EvolvesFromSpecies == nil {
		// Base form or standalone pokemon → Common
		return valueobject.RarityCommon
	}

	// Has a predecessor → at least 2nd stage
	// Check if predecessor also has a predecessor (= 3rd stage)
	predID := extractIDFromURL(species.EvolvesFromSpecies.URL)
	if predID > 0 {
		var predSpecies speciesResponse
		if err := c.fetchURL(ctx, species.EvolvesFromSpecies.URL, &predSpecies); err == nil && predSpecies.EvolvesFromSpecies != nil {
			// 3rd stage evolution → Rare
			return valueobject.RarityRare
		}
	}

	// 2nd stage evolution → Uncommon
	return valueobject.RarityUncommon
}

func (c *Client) fetchURL(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d for %s", resp.StatusCode, url)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// extractIDFromURL extracts the numeric ID from a PokeAPI URL like
// "https://pokeapi.co/api/v2/pokemon-species/1/"
func extractIDFromURL(url string) int {
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return 0
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}
	return id
}
