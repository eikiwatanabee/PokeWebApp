package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TradeStatus string

const (
	TradeStatusOpen      TradeStatus = "open"
	TradeStatusAccepted  TradeStatus = "accepted"
	TradeStatusCancelled TradeStatus = "cancelled"
)

var (
	ErrTradeNotOpen     = errors.New("trade is not open")
	ErrTradeSelfAccept  = errors.New("cannot accept own trade")
	ErrPokemonNotOwned  = errors.New("pokemon not owned by user")
)

type Trade struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	OffererID       uuid.UUID   // User offering the trade
	OfferedPokemonID uuid.UUID  // Pokemon being offered
	RequestedPokemonName string // What they want (e.g., "pikachu") - empty means any
	AccepterID      *uuid.UUID  // User who accepted
	AcceptedPokemonID *uuid.UUID // Pokemon given in exchange
	Status          TradeStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewTrade(tenantID, offererID, offeredPokemonID uuid.UUID, requestedPokemonName string) *Trade {
	now := time.Now()
	return &Trade{
		ID:                   uuid.New(),
		TenantID:             tenantID,
		OffererID:            offererID,
		OfferedPokemonID:     offeredPokemonID,
		RequestedPokemonName: requestedPokemonName,
		Status:               TradeStatusOpen,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

func (t *Trade) Accept(accepterID, acceptedPokemonID uuid.UUID) error {
	if t.Status != TradeStatusOpen {
		return ErrTradeNotOpen
	}
	if t.OffererID == accepterID {
		return ErrTradeSelfAccept
	}
	t.AccepterID = &accepterID
	t.AcceptedPokemonID = &acceptedPokemonID
	t.Status = TradeStatusAccepted
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trade) Cancel(userID uuid.UUID) error {
	if t.Status != TradeStatusOpen {
		return ErrTradeNotOpen
	}
	if t.OffererID != userID {
		return ErrPokemonNotOwned
	}
	t.Status = TradeStatusCancelled
	t.UpdatedAt = time.Now()
	return nil
}
