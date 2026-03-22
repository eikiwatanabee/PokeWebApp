package handler

import (
	"net/http"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type LimitedEventHandler struct {
	getEvents   *query.GetLimitedEventsHandler
	createEvent *command.CreateLimitedEventHandler
}

func NewLimitedEventHandler(
	getEvents *query.GetLimitedEventsHandler,
	createEvent *command.CreateLimitedEventHandler,
) *LimitedEventHandler {
	return &LimitedEventHandler{getEvents: getEvents, createEvent: createEvent}
}

func (h *LimitedEventHandler) GetActiveEvents(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	result, err := h.getEvents.Handle(c.Request.Context(), &query.GetLimitedEventsQuery{
		TenantID: tenantID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *LimitedEventHandler) CreateEvent(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	var body struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Icon        string  `json:"icon"`
		RarityBoost float64 `json:"rarity_boost"`
		StartsAt    string  `json:"starts_at"`
		EndsAt      string  `json:"ends_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startsAt, err := time.Parse(time.RFC3339, body.StartsAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid starts_at format"})
		return
	}
	endsAt, err := time.Parse(time.RFC3339, body.EndsAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ends_at format"})
		return
	}

	eventID, err := h.createEvent.Handle(c.Request.Context(), &command.CreateLimitedEventCommand{
		TenantID:    tenantID,
		Title:       body.Title,
		Description: body.Description,
		Icon:        body.Icon,
		RarityBoost: body.RarityBoost,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event_id": eventID})
}
