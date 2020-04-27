package handlers

import "net/http"

// Ping Checks
func (p *Provider) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}
