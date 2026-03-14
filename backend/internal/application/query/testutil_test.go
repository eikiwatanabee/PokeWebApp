package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// Mock BookRepository
type mockBookRepo struct {
	books map[uuid.UUID]*entity.Book
}

func newMockBookRepo() *mockBookRepo {
	return &mockBookRepo{books: make(map[uuid.UUID]*entity.Book)}
}

func (m *mockBookRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	book, ok := m.books[id]
	if !ok {
		return nil, entity.ErrInvalidTransition
	}
	return book, nil
}

func (m *mockBookRepo) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	return m.FindByID(ctx, id)
}

func (m *mockBookRepo) FindByUserID(ctx context.Context, userID uuid.UUID, status string, tagID string, page, limit int) ([]*entity.Book, int64, error) {
	var result []*entity.Book
	for _, b := range m.books {
		if b.UserID == userID {
			if status != "" && b.Status.String() != status {
				continue
			}
			result = append(result, b)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockBookRepo) Save(ctx context.Context, book *entity.Book) error {
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.books, id)
	return nil
}

// Mock MemoRepository
type mockMemoRepo struct {
	memos map[uuid.UUID]*entity.Memo
}

func newMockMemoRepo() *mockMemoRepo {
	return &mockMemoRepo{memos: make(map[uuid.UUID]*entity.Memo)}
}

func (m *mockMemoRepo) FindByBookID(ctx context.Context, bookID uuid.UUID) ([]*entity.Memo, error) {
	var result []*entity.Memo
	for _, memo := range m.memos {
		if memo.BookID == bookID {
			result = append(result, memo)
		}
	}
	return result, nil
}

func (m *mockMemoRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Memo, error) {
	memo, ok := m.memos[id]
	if !ok {
		return nil, entity.ErrEmptyMemoContent
	}
	return memo, nil
}

func (m *mockMemoRepo) Save(ctx context.Context, memo *entity.Memo) error {
	m.memos[memo.ID] = memo
	return nil
}

func (m *mockMemoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.memos, id)
	return nil
}

// Mock TagRepository
type mockTagRepo struct {
	tags map[uuid.UUID]*entity.Tag
}

func newMockTagRepo() *mockTagRepo {
	return &mockTagRepo{tags: make(map[uuid.UUID]*entity.Tag)}
}

func (m *mockTagRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	tag, ok := m.tags[id]
	if !ok {
		return nil, entity.ErrEmptyTagName
	}
	return tag, nil
}

func (m *mockTagRepo) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.Tag, error) {
	var result []*entity.Tag
	for _, t := range m.tags {
		if t.TenantID == tenantID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTagRepo) FindByName(ctx context.Context, tenantID uuid.UUID, name string) (*entity.Tag, error) {
	for _, t := range m.tags {
		if t.TenantID == tenantID && t.Name == name {
			return t, nil
		}
	}
	return nil, entity.ErrEmptyTagName
}

func (m *mockTagRepo) Save(ctx context.Context, tag *entity.Tag) error {
	m.tags[tag.ID] = tag
	return nil
}

func (m *mockTagRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.tags, id)
	return nil
}

// Mock PokemonRepository
type mockPokemonRepo struct {
	pokemons map[uuid.UUID]*entity.UserPokemon
}

func newMockPokemonRepo() *mockPokemonRepo {
	return &mockPokemonRepo{pokemons: make(map[uuid.UUID]*entity.UserPokemon)}
}

func (m *mockPokemonRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserPokemon, error) {
	var result []*entity.UserPokemon
	for _, p := range m.pokemons {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockPokemonRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.UserPokemon, error) {
	p, ok := m.pokemons[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockPokemonRepo) Save(ctx context.Context, pokemon *entity.UserPokemon) error {
	m.pokemons[pokemon.ID] = pokemon
	return nil
}

func newTestPokemonInfo() valueobject.PokemonInfo {
	return valueobject.PokemonInfo{
		PokemonID: 25,
		Name:      "pikachu",
		SpriteURL: "https://example.com/pikachu.png",
		Types:     []string{"electric"},
	}
}
