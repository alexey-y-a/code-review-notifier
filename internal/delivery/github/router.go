package github

import (
	"net/http"
)

func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/github/webhook", handler.HandleWebhook)
	return mux
}
