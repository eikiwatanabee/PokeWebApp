package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type WebhookVerifier struct {
	secret []byte
}

func NewWebhookVerifier(secret string) *WebhookVerifier {
	return &WebhookVerifier{secret: []byte(secret)}
}

func (v *WebhookVerifier) Verify(payload []byte, signature string) bool {
	if v.secret == nil || len(v.secret) == 0 {
		// No secret configured, skip verification (dev mode)
		return true
	}

	sig := strings.TrimPrefix(signature, "sha256=")
	mac := hmac.New(sha256.New, v.secret)
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sig), []byte(expected))
}

func (v *WebhookVerifier) Sign(payload []byte) string {
	mac := hmac.New(sha256.New, v.secret)
	mac.Write(payload)
	return fmt.Sprintf("sha256=%s", hex.EncodeToString(mac.Sum(nil)))
}
