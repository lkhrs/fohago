package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lkhrs/fohago/antispam"
	"github.com/writefreely/go-akismet"
)

type Check struct{}

func (c *Check) blocklist(sub FormSubmission, fh FormHandler) (bool, error) {
	global := fh.Config.Global.Blocklist
	form := sub.FormCfg.Blocklist
	combined := append(global, form...)
	for _, term := range combined {
		if strings.Contains(sub.Body[sub.FormCfg.Fields.Message], term) {
			return false, fmt.Errorf("message contains blocklist term \"%v\"", term)
		}
	}
	return true, nil
}

func (c *Check) honeypot(sub FormSubmission) (bool, error) {
	if field, exists := sub.Body[sub.FormCfg.Fields.Honeypot]; exists {
		if field != "" {
			return false, errors.New("honeypot field is not empty")
		}
	}
	return true, nil
}

func (c *Check) turnstile(sub FormSubmission) (bool, error) {
	if sub.FormCfg.TurnstileKey == "" {
		return true, nil
	}
	secret := sub.FormCfg.TurnstileKey
	token := sub.Body["cf-turnstile-response"]
	return antispam.Turnstile(secret, token)
}

func (c *Check) akismet(sub FormSubmission, fh FormHandler, isTest bool, userRole string) (bool, error) {
	if fh.Config.Api.Akismet == "" {
		return true, errors.New("no Akismet key provided")
	}
	isSpam, err := akismet.Check(&akismet.Comment{
		Blog:               fh.Config.Global.BaseUrl, // required
		UserIP:             sub.UserIP,               // required
		UserAgent:          sub.UserAgent,            // required
		Referrer:           sub.Referrer,
		CommentType:        "contact‑form",
		CommentAuthor:      sub.Body[sub.FormCfg.Fields.Name],
		CommentAuthorEmail: sub.Body[sub.FormCfg.Fields.Email],
		CommentContent:     sub.Body[sub.FormCfg.Fields.Message],
		CommentDate:        time.Now(),
		UserRole:           userRole,
		Test:               isTest,
	}, fh.Config.Api.Akismet)
	return isSpam, err
}

// checkSpam checks the form submission for spam
// returns true if the checks pass, false if spam is detected
func (fh *FormHandler) checkSpam(sub FormSubmission) bool {
	check := &Check{}
	if pass, err := check.honeypot(sub); !pass {
		log.Println("Honeypot check failed:", err)
		return false
	}
	if pass, err := check.blocklist(sub, *fh); !pass {
		log.Println(err)
		return false
	}
	if pass, err := check.turnstile(sub); !pass {
		log.Println("Turnstile check failed:", err)
		return false
	}
	if isSpam, err := check.akismet(sub, *fh, false, ""); isSpam {
		log.Println("Akismet check failed:", err)
		return false
	}
	return true
}
