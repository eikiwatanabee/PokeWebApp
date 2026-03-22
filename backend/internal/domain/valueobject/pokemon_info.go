package valueobject

type PokemonRarity string

const (
	RarityCommon    PokemonRarity = "common"
	RarityUncommon  PokemonRarity = "uncommon"
	RarityRare      PokemonRarity = "rare"
	RarityEpic      PokemonRarity = "epic"
	RarityLegendary PokemonRarity = "legendary"
)

type PokemonInfo struct {
	PokemonID int           `json:"pokemon_id"`
	Name      string        `json:"name"`
	SpriteURL string        `json:"sprite_url"`
	Types     []string      `json:"types"`
	Rarity    PokemonRarity `json:"rarity"`
}
