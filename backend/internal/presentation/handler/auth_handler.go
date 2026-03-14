package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	googleOAuth *auth.GoogleOAuth
	jwtManager  *auth.JWTManager
	userRepo    repository.UserRepository
	tenantRepo  repository.TenantRepository
}

func NewAuthHandler(
	googleOAuth *auth.GoogleOAuth,
	jwtManager *auth.JWTManager,
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
) *AuthHandler {
	return &AuthHandler{
		googleOAuth: googleOAuth,
		jwtManager:  jwtManager,
		userRepo:    userRepo,
		tenantRepo:  tenantRepo,
	}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	url := h.googleOAuth.GetAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	token, err := h.googleOAuth.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	userInfo, err := h.googleOAuth.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// Find or create user
	user, err := h.userRepo.FindByGoogleID(c.Request.Context(), userInfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if user == nil {
		// Create default tenant and user
		tenant, err := entity.NewTenant(userInfo.Name + "'s Team")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant"})
			return
		}
		if err := h.tenantRepo.Save(c.Request.Context(), tenant); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save tenant"})
			return
		}

		user, err = entity.NewUser(tenant.ID, userInfo.ID, userInfo.Email, userInfo.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		user.PromoteToAdmin() // First user is admin
		if err := h.userRepo.Save(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
			return
		}
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.TenantID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.ID.String(),
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, err := h.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.TenantID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
