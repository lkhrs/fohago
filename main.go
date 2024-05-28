package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/lkhrs/fohago/middleware"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	fh := NewFormHandler(config)

	mux.HandleFunc("POST /{id}", fh.handleFormSubmission)
	mux.HandleFunc("GET /test.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test.html")
	})

	handler := middleware.Logging(mux)
	handler = middleware.PanicRecovery(handler)

	http.ListenAndServe(":"+strconv.Itoa(config.Global.Port), handler)
}
