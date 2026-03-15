package valueobject

type PokemonInfo struct {
	PokemonID   int      `json:"pokemon_id"`
	Name        string   `json:"name"`
	SpriteURL   string   `json:"sprite_url"`
	Types       []string `json:"types"`
}
