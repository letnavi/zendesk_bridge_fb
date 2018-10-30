package main

import (
	"net/http"
	"os"
)

type Feed struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Period string `json:"period"`
	Values []struct {
		Value int `json:"value"`
	} `json:"values"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Verification Endpoint
func Verify(t string, w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	token := r.URL.Query().Get("hub.verify_token")

	if token == os.Getenv(t) {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}
