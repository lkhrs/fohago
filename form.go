package main

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type FormHandler struct {
	Config         *Config
	FormSubmission FormSubmission
}

type FormSubmission struct {
	Id        string
	Body      FormBody
	FormCfg   FormConfig
	UserAgent string
	UserIP    string
	Referrer  string
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
	successRedirect := fh.Config.Global.BaseUrl + "/success.html"
	if submission.FormCfg.Redirects.Success != "" {
		successRedirect = submission.FormCfg.Redirects.Success
	}
	http.Redirect(w, r, successRedirect, http.StatusFound)
}

func (fh *FormHandler) getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
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

	fields := make(FormBody)
	p := bluemonday.StrictPolicy()
	for k, v := range r.Form {
		fields[k] = p.Sanitize(v[0])
	}

	submission := FormSubmission{
		Id:        id,
		Body:      fields,
		FormCfg:   formCfg,
		UserAgent: r.UserAgent(),
		UserIP:    fh.getClientIP(r),
		Referrer:  r.Referer(),
	}

	return submission
}

// sendMail sends the form submission as an email
// returns true if the email was sent successfully, false otherwise
func (fh *FormHandler) sendMail(sub FormSubmission) bool {
	return buildAndSend(fh.Config, sub) == nil
}
