package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alexey-y-a/code-review-notifier/internal/model"
	"github.com/alexey-y-a/code-review-notifier/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	service service.Service
	secret  string
}

func NewHandler(service service.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if !h.validateSignature(r, signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("Failed to read request body", zap.Error(err))
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	var payload struct {
		Action      string `json:"action"`
		PullRequest struct {
			Title    string `json:"title"`
			HTMLURL  string `json:"html_url"`
			Assignee struct {
				Login string `json:"login"`
			} `json:"assignee"`
		} `json:"pull_request"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		zap.L().Error("Failed to parse webhook payload", zap.Error(err))
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if payload.Action != "assigned" || payload.PullRequest.Assignee.Login == "" {
		zap.L().Debug("Ignoring event", zap.String("action", payload.Action))
		w.WriteHeader(http.StatusOK)
		return
	}

	event := &model.PullRequestEvent{
		Action:     payload.Action,
		Assignee:   payload.PullRequest.Assignee.Login,
		Title:      payload.PullRequest.Title,
		HTMLURL:    payload.PullRequest.HTMLURL,
		Repository: payload.Repository.FullName,
	}

	if err := h.service.HandleGitHubEvent(event); err != nil {
		zap.L().Error("Failed to handle GitHub event", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) validateSignature(r *http.Request, signature string) bool {
	if h.secret == "" {
		return true
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("Failed to read body for signature validation", zap.Error(err))
		return false
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	hash := hmac.New(sha256.New, []byte(h.secret))
	hash.Write(body)
	expected := "sha256=" + hex.EncodeToString(hash.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
