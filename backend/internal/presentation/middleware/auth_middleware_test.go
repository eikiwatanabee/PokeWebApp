package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testJWTSecret = "test-secret-key-for-jwt"

func TestAuthRequired_MissingAuthorizationHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/books", nil)

	jwtManager := auth.NewJWTManager(testJWTSecret)
	handler := AuthRequired(jwtManager)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "missing authorization header" {
		t.Errorf("expected 'missing authorization header', got %v", resp["error"])
	}
}

func TestAuthRequired_InvalidFormat_NoBearerPrefix(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/books", nil)
	c.Request.Header.Set("Authorization", "Token some-token")

	jwtManager := auth.NewJWTManager(testJWTSecret)
	handler := AuthRequired(jwtManager)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "invalid authorization format" {
		t.Errorf("expected 'invalid authorization format', got %v", resp["error"])
	}
}

func TestAuthRequired_InvalidFormat_BearerOnly(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/books", nil)
	c.Request.Header.Set("Authorization", "Bearer")

	jwtManager := auth.NewJWTManager(testJWTSecret)
	handler := AuthRequired(jwtManager)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "invalid authorization format" {
		t.Errorf("expected 'invalid authorization format', got %v", resp["error"])
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/books", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	jwtManager := auth.NewJWTManager(testJWTSecret)
	handler := AuthRequired(jwtManager)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "invalid or expired token" {
		t.Errorf("expected 'invalid or expired token', got %v", resp["error"])
	}
}

func TestAuthRequired_ValidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager(testJWTSecret)
	userID := uuid.New()
	tenantID := uuid.New()
	token, err := jwtManager.GenerateAccessToken(userID, tenantID, "test@example.com", "user")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/books", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := AuthRequired(jwtManager)
	handler(c)

	// Should not abort - status stays at 200 (default)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if c.GetString("user_id") != userID.String() {
		t.Errorf("expected user_id %s, got %s", userID.String(), c.GetString("user_id"))
	}
	if c.GetString("tenant_id") != tenantID.String() {
		t.Errorf("expected tenant_id %s, got %s", tenantID.String(), c.GetString("tenant_id"))
	}
	if c.GetString("email") != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", c.GetString("email"))
	}
	if c.GetString("role") != "user" {
		t.Errorf("expected role user, got %s", c.GetString("role"))
	}
}

func TestAdminRequired_NoRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/admin", nil)

	handler := AdminRequired()
	handler(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "admin access required" {
		t.Errorf("expected 'admin access required', got %v", resp["error"])
	}
}

func TestAdminRequired_UserRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/admin", nil)
	c.Set("role", "user")

	handler := AdminRequired()
	handler(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "admin access required" {
		t.Errorf("expected 'admin access required', got %v", resp["error"])
	}
}

func TestAdminRequired_AdminRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/admin", nil)
	c.Set("role", "admin")

	handler := AdminRequired()
	handler(c)

	// Should not abort - status stays at 200 (default)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
