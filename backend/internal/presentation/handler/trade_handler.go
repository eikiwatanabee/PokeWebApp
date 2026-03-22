package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type TradeHandler struct {
	createTrade *command.CreateTradeHandler
	acceptTrade *command.AcceptTradeHandler
	cancelTrade *command.CancelTradeHandler
	getTrades   *query.GetTradesHandler
}

func NewTradeHandler(
	createTrade *command.CreateTradeHandler,
	acceptTrade *command.AcceptTradeHandler,
	cancelTrade *command.CancelTradeHandler,
	getTrades *query.GetTradesHandler,
) *TradeHandler {
	return &TradeHandler{
		createTrade: createTrade,
		acceptTrade: acceptTrade,
		cancelTrade: cancelTrade,
		getTrades:   getTrades,
	}
}

func (h *TradeHandler) GetTrades(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	result, err := h.getTrades.Handle(c.Request.Context(), &query.GetTradesQuery{
		TenantID: tenantID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TradeHandler) CreateTrade(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	var body struct {
		OfferedPokemonID     string `json:"offered_pokemon_id"`
		RequestedPokemonName string `json:"requested_pokemon_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tradeID, err := h.createTrade.Handle(c.Request.Context(), &command.CreateTradeCommand{
		TenantID:             tenantID,
		UserID:               userID,
		OfferedPokemonID:     body.OfferedPokemonID,
		RequestedPokemonName: body.RequestedPokemonName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"trade_id": tradeID})
}

func (h *TradeHandler) AcceptTrade(c *gin.Context) {
	userID := c.GetString("user_id")
	tradeID := c.Param("id")

	var body struct {
		OfferedPokemonID string `json:"offered_pokemon_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.acceptTrade.Handle(c.Request.Context(), &command.AcceptTradeCommand{
		TradeID:          tradeID,
		UserID:           userID,
		OfferedPokemonID: body.OfferedPokemonID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trade accepted"})
}

func (h *TradeHandler) CancelTrade(c *gin.Context) {
	userID := c.GetString("user_id")
	tradeID := c.Param("id")

	err := h.cancelTrade.Handle(c.Request.Context(), &command.CancelTradeCommand{
		TradeID: tradeID,
		UserID:  userID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trade cancelled"})
}
