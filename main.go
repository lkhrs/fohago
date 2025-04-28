package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/lkhrs/fohago/middleware"
)

func main() {
	// Set up logging
	serviceLogger := slog.New(ServiceLogHandler())
	slog.SetDefault(serviceLogger)
	accessLogger := slog.New(AccessLogHandler())

	// Load config
	config := loadConfig("fohago.toml")

	// Set up HTTP handler
	mux := http.NewServeMux()
	fh := NewFormHandler(config)

	// Routes
	mux.HandleFunc("POST /{id}", fh.handleFormSubmission)
	mux.HandleFunc("GET /test.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test.html")
	})
	mux.HandleFunc("GET /success.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./success.html")
	})

	// Middleware
	handler := middleware.Logging(mux, accessLogger)
	handler = middleware.PanicRecovery(handler)

	// Start server
	http.ListenAndServe(":"+strconv.Itoa(config.Global.Port), handler)
}
