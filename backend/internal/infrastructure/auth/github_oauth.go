package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GitHubOAuth struct {
	config *oauth2.Config
}

func NewGitHubOAuth(clientID, clientSecret, redirectURL string) *GitHubOAuth {
	return &GitHubOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (g *GitHubOAuth) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GitHubOAuth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

func (g *GitHubOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GitHubUserInfo, error) {
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	// If email is empty, fetch from emails endpoint
	if userInfo.Email == "" {
		email, err := g.fetchPrimaryEmail(ctx, token)
		if err == nil && email != "" {
			userInfo.Email = email
		}
	}

	if userInfo.Name == "" {
		userInfo.Name = userInfo.Login
	}

	return &userInfo, nil
}

func (g *GitHubOAuth) fetchPrimaryEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary {
			return e.Email, nil
		}
	}
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	return "", nil
}
