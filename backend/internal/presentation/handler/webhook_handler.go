package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	verifier           *auth.WebhookVerifier
	processGitHubEvent *command.ProcessGitHubEventHandler
}

func NewWebhookHandler(
	verifier *auth.WebhookVerifier,
	processGitHubEvent *command.ProcessGitHubEventHandler,
) *WebhookHandler {
	return &WebhookHandler{
		verifier:           verifier,
		processGitHubEvent: processGitHubEvent,
	}
}

func (h *WebhookHandler) HandleGitHubWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Verify webhook signature
	signature := c.GetHeader("X-Hub-Signature-256")
	if !h.verifier.Verify(body, signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	eventType := c.GetHeader("X-GitHub-Event")
	ctx := c.Request.Context()

	switch eventType {
	case "push":
		h.handlePush(ctx, body, c)
	case "pull_request":
		h.handlePullRequest(ctx, body, c)
	case "pull_request_review":
		h.handlePullRequestReview(ctx, body, c)
	case "issues":
		h.handleIssues(ctx, body, c)
	default:
		c.JSON(http.StatusOK, gin.H{"message": "event type not tracked", "event": eventType})
	}
}

type pushPayload struct {
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Username string `json:"username"`
		} `json:"author"`
	} `json:"commits"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (h *WebhookHandler) handlePush(ctx interface{ Done() <-chan struct{} }, body []byte, c *gin.Context) {
	var payload pushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid push payload"})
		return
	}

	processed := 0
	for _, commit := range payload.Commits {
		if commit.Author.Username == "" {
			continue
		}
		title := commit.Message
		if len(title) > 100 {
			title = title[:100]
		}
		result, err := h.processGitHubEvent.Handle(c.Request.Context(), &command.ProcessGitHubEventCommand{
			GitHubUsername: commit.Author.Username,
			EventType:     entity.EventCommit,
			RepoName:      payload.Repository.FullName,
			Title:         title,
			URL:           commit.URL,
		})
		if err != nil {
			log.Printf("error processing commit %s: %v", commit.ID, err)
			continue
		}
		if result != nil {
			processed++
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "push processed", "commits_tracked": processed})
}

type prPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Title    string `json:"title"`
		HTMLURL  string `json:"html_url"`
		Merged   bool   `json:"merged"`
		User     struct {
			Login string `json:"login"`
		} `json:"user"`
		MergedBy *struct {
			Login string `json:"login"`
		} `json:"merged_by"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (h *WebhookHandler) handlePullRequest(ctx interface{ Done() <-chan struct{} }, body []byte, c *gin.Context) {
	var payload prPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid PR payload"})
		return
	}

	switch payload.Action {
	case "opened":
		_, err := h.processGitHubEvent.Handle(c.Request.Context(), &command.ProcessGitHubEventCommand{
			GitHubUsername: payload.PullRequest.User.Login,
			EventType:     entity.EventPROpen,
			RepoName:      payload.Repository.FullName,
			Title:         payload.PullRequest.Title,
			URL:           payload.PullRequest.HTMLURL,
		})
		if err != nil {
			log.Printf("error processing PR open: %v", err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "pr_open processed"})

	case "closed":
		if payload.PullRequest.Merged {
			// Give XP to the PR author for the merge
			username := payload.PullRequest.User.Login
			result, err := h.processGitHubEvent.Handle(c.Request.Context(), &command.ProcessGitHubEventCommand{
				GitHubUsername: username,
				EventType:     entity.EventPRMerge,
				RepoName:      payload.Repository.FullName,
				Title:         payload.PullRequest.Title,
				URL:           payload.PullRequest.HTMLURL,
			})
			if err != nil {
				log.Printf("error processing PR merge: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process merge"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "pr_merge processed", "result": result})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "pr closed without merge, skipped"})
		}

	default:
		c.JSON(http.StatusOK, gin.H{"message": "pr action not tracked", "action": payload.Action})
	}
}

type reviewPayload struct {
	Action string `json:"action"`
	Review struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		HTMLURL string `json:"html_url"`
	} `json:"review"`
	PullRequest struct {
		Title string `json:"title"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (h *WebhookHandler) handlePullRequestReview(ctx interface{ Done() <-chan struct{} }, body []byte, c *gin.Context) {
	var payload reviewPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review payload"})
		return
	}

	if payload.Action != "submitted" {
		c.JSON(http.StatusOK, gin.H{"message": "review action not tracked"})
		return
	}

	_, err := h.processGitHubEvent.Handle(c.Request.Context(), &command.ProcessGitHubEventCommand{
		GitHubUsername: payload.Review.User.Login,
		EventType:     entity.EventReview,
		RepoName:      payload.Repository.FullName,
		Title:         "Review: " + payload.PullRequest.Title,
		URL:           payload.Review.HTMLURL,
	})
	if err != nil {
		log.Printf("error processing review: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "review processed"})
}

type issuePayload struct {
	Action string `json:"action"`
	Issue  struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"issue"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (h *WebhookHandler) handleIssues(ctx interface{ Done() <-chan struct{} }, body []byte, c *gin.Context) {
	var payload issuePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid issue payload"})
		return
	}

	if payload.Action != "closed" {
		c.JSON(http.StatusOK, gin.H{"message": "issue action not tracked"})
		return
	}

	_, err := h.processGitHubEvent.Handle(c.Request.Context(), &command.ProcessGitHubEventCommand{
		GitHubUsername: payload.Issue.User.Login,
		EventType:     entity.EventIssueClose,
		RepoName:      payload.Repository.FullName,
		Title:         payload.Issue.Title,
		URL:           payload.Issue.HTMLURL,
	})
	if err != nil {
		log.Printf("error processing issue close: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "issue_close processed"})
}
