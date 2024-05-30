package main

import (
	"errors"

	"github.com/lkhrs/fohago/antispam"
)

type Check struct{}

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
