package main

import (
	"github.com/microcosm-cc/bluemonday"
	"log"
	"net/http"
)

type FormHandler struct {
	FormSubmission FormSubmission
	Config         *Config
}

type FormSubmission struct {
	Id      string
	Fields  FormFields
	FormCfg FormConfig
}

func NewFormHandler(conf *Config) *FormHandler {
	fh := &FormHandler{Config: conf}
	return fh
}

func (fh *FormHandler) handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	submission := fh.process(w, r)
	if !fh.checkSpam(submission) {
		http.Error(w, "Spam detected", http.StatusBadRequest)
		return
	}
	success := fh.sendMail(submission)
	if !success {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

// process parses the form submission and returns a FormSubmission struct
func (fh *FormHandler) process(w http.ResponseWriter, r *http.Request) FormSubmission {
	id := r.URL.Path[1:]
	formCfg, exists := fh.Config.Forms[id]
	if !exists {
		http.NotFound(w, r)
	}

	if err := r.ParseForm(); err != nil {
		log.Println("Failed to parse form:", err)
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
	}

	fields := make(FormFields)
	p := bluemonday.StrictPolicy()
	for k, v := range r.Form {
		fields[k] = p.Sanitize(v[0])
	}

	submission := FormSubmission{
		Id:      id,
		Fields:  fields,
		FormCfg: formCfg,
	}

	return submission
}

// checkSpam checks the form submission for spam
// returns true if the checks pass, false if spam is detected
func (fh *FormHandler) checkSpam(sub FormSubmission) bool {
	check := &Check{}
	if pass, err := check.honeypot(sub); !pass {
		log.Println("Honeypot check failed:", err)
		return false
	}
	if pass, err := check.turnstile(sub); !pass {
		log.Println("Turnstile check failed:", err)
		return false
	}
	return true
}

// sendMail sends the form submission as an email
// returns true if the email was sent successfully, false otherwise
func (fh *FormHandler) sendMail(sub FormSubmission) bool {
	return buildAndSend(fh.Config, sub) == nil
}
