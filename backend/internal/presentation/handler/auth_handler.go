package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	googleOAuth *auth.GoogleOAuth
	githubOAuth *auth.GitHubOAuth
	jwtManager  *auth.JWTManager
	userRepo    repository.UserRepository
	tenantRepo  repository.TenantRepository
}

func NewAuthHandler(
	googleOAuth *auth.GoogleOAuth,
	githubOAuth *auth.GitHubOAuth,
	jwtManager *auth.JWTManager,
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
) *AuthHandler {
	return &AuthHandler{
		googleOAuth: googleOAuth,
		githubOAuth: githubOAuth,
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

	h.respondWithTokens(c, user)
}

// GitHubLogin returns the GitHub OAuth URL
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	url := h.githubOAuth.GetAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GitHubCallback handles the GitHub OAuth callback
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	token, err := h.githubOAuth.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	userInfo, err := h.githubOAuth.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	githubID := fmt.Sprintf("%d", userInfo.ID)
	googleID := "github:" + githubID

	// Find existing user by GoogleID (which stores github: prefix for GitHub users)
	user, err := h.userRepo.FindByGoogleID(c.Request.Context(), googleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if user == nil {
		// Create default tenant and user
		tenant, err := entity.NewTenant(userInfo.Login + "'s Team")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant"})
			return
		}
		if err := h.tenantRepo.Save(c.Request.Context(), tenant); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save tenant"})
			return
		}

		email := userInfo.Email
		if email == "" {
			email = userInfo.Login + "@github.com"
		}

		user, err = entity.NewGitHubUser(tenant.ID, githubID, userInfo.Login, email, userInfo.Name, userInfo.AvatarURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		user.PromoteToAdmin()
		if err := h.userRepo.Save(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
			return
		}
	}

	h.respondWithTokens(c, user)
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

// DevLogin creates a dummy user for local development (no OAuth required).
func (h *AuthHandler) DevLogin(c *gin.Context) {
	ctx := c.Request.Context()

	devGoogleID := "dev-user-local"
	user, err := h.userRepo.FindByGoogleID(ctx, devGoogleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if user == nil {
		tenant, err := entity.NewTenant("Dev Team")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant"})
			return
		}
		if err := h.tenantRepo.Save(ctx, tenant); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save tenant"})
			return
		}

		user, err = entity.NewUser(tenant.ID, devGoogleID, "dev@localhost", "Dev User")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		user.GitHubUsername = "dev-user"
		user.GitHubID = "dev-user"
		user.PromoteToAdmin()
		if err := h.userRepo.Save(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
			return
		}
	}

	h.respondWithTokens(c, user)
}

func (h *AuthHandler) respondWithTokens(c *gin.Context, user *entity.User) {
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
			"id":              user.ID.String(),
			"name":            user.Name,
			"email":           user.Email,
			"role":            user.Role,
			"github_username": user.GitHubUsername,
			"avatar_url":      user.AvatarURL,
			"level":           user.Level,
			"total_xp":        user.TotalXP,
		},
	})
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
