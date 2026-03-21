package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testWebhookSecret = "test-secret"

func setupWebhookTest(eventType string, payload interface{}) (*httptest.ResponseRecorder, *gin.Context, []byte) {
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", eventType)

	verifier := auth.NewWebhookVerifier(testWebhookSecret)
	sig := verifier.Sign(body)
	req.Header.Set("X-Hub-Signature-256", sig)

	c.Request = req
	return w, c, body
}

func TestWebhookHandler_InvalidSignature(t *testing.T) {
	body := []byte(`{}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-Hub-Signature-256", "sha256=invalidsig")
	c.Request = req

	verifier := auth.NewWebhookVerifier(testWebhookSecret)
	h := NewWebhookHandler(verifier, nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestWebhookHandler_NoSecret_SkipsVerification(t *testing.T) {
	body := []byte(`{}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "unknown_event")
	c.Request = req

	verifier := auth.NewWebhookVerifier("") // empty = dev mode
	h := NewWebhookHandler(verifier, nil)
	h.HandleGitHubWebhook(c)

	// Should not get 401, but process as unknown event
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for unknown event, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "event type not tracked" {
		t.Errorf("expected 'event type not tracked', got %v", resp["message"])
	}
}

func TestWebhookHandler_UnknownEventType(t *testing.T) {
	w, c, _ := setupWebhookTest("workflow_run", map[string]string{})

	h := NewWebhookHandler(auth.NewWebhookVerifier(testWebhookSecret), nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "event type not tracked" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
}

func TestWebhookHandler_DeploymentNoUsername(t *testing.T) {
	w, c, _ := setupWebhookTest("deployment", map[string]interface{}{
		"deployment": map[string]interface{}{
			"creator":     map[string]string{},
			"environment": "production",
		},
		"sender":     map[string]string{},
		"repository": map[string]string{"full_name": "test/repo"},
	})

	h := NewWebhookHandler(auth.NewWebhookVerifier(testWebhookSecret), nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "no deployer username found" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
}

func TestWebhookHandler_PushInvalidPayload(t *testing.T) {
	body := []byte(`not json`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "push")

	verifier := auth.NewWebhookVerifier("")
	req.Header.Set("X-Hub-Signature-256", verifier.Sign(body))
	c.Request = req

	h := NewWebhookHandler(verifier, nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid push payload, got %d", w.Code)
	}
}

func TestWebhookHandler_PRInvalidPayload(t *testing.T) {
	body := []byte(`{invalid}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "pull_request")

	verifier := auth.NewWebhookVerifier("")
	req.Header.Set("X-Hub-Signature-256", verifier.Sign(body))
	c.Request = req

	h := NewWebhookHandler(verifier, nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_ReviewInvalidPayload(t *testing.T) {
	body := []byte(`{invalid}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "pull_request_review")

	verifier := auth.NewWebhookVerifier("")
	req.Header.Set("X-Hub-Signature-256", verifier.Sign(body))
	c.Request = req

	h := NewWebhookHandler(verifier, nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_IssuesInvalidPayload(t *testing.T) {
	body := []byte(`{invalid}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewReader(body))
	req.Header.Set("X-GitHub-Event", "issues")

	verifier := auth.NewWebhookVerifier("")
	req.Header.Set("X-Hub-Signature-256", verifier.Sign(body))
	c.Request = req

	h := NewWebhookHandler(verifier, nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhookHandler_IssuesNonClosedAction(t *testing.T) {
	payload := map[string]interface{}{
		"action": "opened",
		"issue": map[string]interface{}{
			"title":    "test issue",
			"html_url": "https://github.com/test/repo/issues/1",
			"user":     map[string]string{"login": "testuser"},
		},
		"repository": map[string]string{"full_name": "test/repo"},
	}

	w, c, _ := setupWebhookTest("issues", payload)

	h := NewWebhookHandler(auth.NewWebhookVerifier(testWebhookSecret), nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "issue action not tracked" {
		t.Errorf("expected 'issue action not tracked', got %v", resp["message"])
	}
}

func TestWebhookHandler_ReviewNonSubmittedAction(t *testing.T) {
	payload := map[string]interface{}{
		"action": "dismissed",
		"review": map[string]interface{}{
			"user":     map[string]string{"login": "reviewer"},
			"html_url": "https://github.com/test/repo/pulls/1#pullrequestreview-1",
		},
		"pull_request": map[string]string{"title": "test PR"},
		"repository":   map[string]string{"full_name": "test/repo"},
	}

	w, c, _ := setupWebhookTest("pull_request_review", payload)

	h := NewWebhookHandler(auth.NewWebhookVerifier(testWebhookSecret), nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "review action not tracked" {
		t.Errorf("expected 'review action not tracked', got %v", resp["message"])
	}
}

func TestWebhookHandler_PRClosedWithoutMerge(t *testing.T) {
	payload := map[string]interface{}{
		"action": "closed",
		"pull_request": map[string]interface{}{
			"title":    "test PR",
			"html_url": "https://github.com/test/repo/pull/1",
			"merged":   false,
			"user":     map[string]string{"login": "author"},
		},
		"repository": map[string]string{"full_name": "test/repo"},
	}

	w, c, _ := setupWebhookTest("pull_request", payload)

	h := NewWebhookHandler(auth.NewWebhookVerifier(testWebhookSecret), nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "pr closed without merge, skipped" {
		t.Errorf("expected 'pr closed without merge, skipped', got %v", resp["message"])
	}
}

func TestWebhookHandler_PRUntrackedAction(t *testing.T) {
	payload := map[string]interface{}{
		"action": "synchronize",
		"pull_request": map[string]interface{}{
			"title":    "test PR",
			"html_url": "https://github.com/test/repo/pull/1",
			"user":     map[string]string{"login": "author"},
		},
		"repository": map[string]string{"full_name": "test/repo"},
	}

	w, c, _ := setupWebhookTest("pull_request", payload)

	h := NewWebhookHandler(auth.NewWebhookVerifier(testWebhookSecret), nil)
	h.HandleGitHubWebhook(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "pr action not tracked" {
		t.Errorf("expected 'pr action not tracked', got %v", resp["message"])
	}
}
