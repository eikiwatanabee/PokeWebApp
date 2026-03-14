---
name: add-command
description: Generate a CQRS Command with handler for write operations. Use when adding new write/mutation functionality.
argument-hint: [CommandName]
---

# CQRS Command Generator

Generate a new Command (write side of CQRS) for the PokeBookManager project.

## Command Name: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for domain context
2. Generate the following files:

### Command DTO (`backend/internal/application/command/$ARGUMENTS.go`)
- Input struct with validation tags
- Handler struct with dependencies (repositories, UoW)
- `Handle(ctx context.Context, cmd *{CommandName}Command) (*{Result}, error)` method

### Pattern

```go
package command

import (
    "context"
    "github.com/your-org/pokebookmanager/internal/application/uow"
    "github.com/your-org/pokebookmanager/internal/domain/repository"
)

type FinishReadingCommand struct {
    BookID string
    UserID string
}

type FinishReadingResult struct {
    PokemonID   int
    PokemonName string
}

type FinishReadingHandler struct {
    uow         uow.UnitOfWork
    bookRepo    repository.BookRepository
    pokemonRepo repository.PokemonRepository
    pokemonSvc  service.PokemonGachaService
}

func NewFinishReadingHandler(
    uow uow.UnitOfWork,
    bookRepo repository.BookRepository,
    pokemonRepo repository.PokemonRepository,
    pokemonSvc service.PokemonGachaService,
) *FinishReadingHandler {
    return &FinishReadingHandler{
        uow: uow, bookRepo: bookRepo,
        pokemonRepo: pokemonRepo, pokemonSvc: pokemonSvc,
    }
}

func (h *FinishReadingHandler) Handle(ctx context.Context, cmd *FinishReadingCommand) (*FinishReadingResult, error) {
    var result *FinishReadingResult
    err := h.uow.Do(ctx, func(ctx context.Context) error {
        // 1. Fetch book (lock order: Book=3)
        book, err := h.bookRepo.FindByIDForUpdate(ctx, cmd.BookID)
        if err != nil {
            return err
        }
        // 2. Domain logic
        if err := book.Finish(); err != nil {
            return err
        }
        // 3. Save book
        if err := h.bookRepo.Save(ctx, book); err != nil {
            return err
        }
        // 4. Pokemon gacha (lock order: UserPokemon=6, after Book=3 ✓)
        pokemon, err := h.pokemonSvc.Draw(ctx)
        if err != nil {
            return err
        }
        // 5. Save pokemon
        if err := h.pokemonRepo.Save(ctx, pokemon); err != nil {
            return err
        }
        result = &FinishReadingResult{
            PokemonID: pokemon.PokemonID, PokemonName: pokemon.Name,
        }
        return nil
    })
    return result, err
}
```

## Rules
- Commands are write operations only
- Always use UoW for transactional consistency
- Respect lock ordering: Tenant(1) → User(2) → Book(3) → Memo(4) → Tag(5) → UserPokemon(6)
- Return minimal result (ID + essential fields)
- Domain logic stays in entities, not in handlers
