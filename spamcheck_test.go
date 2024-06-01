package main

import (
	"errors"
	"testing"
)

func TestCheck_blocklist(t *testing.T) {
	c := &Check{}
	fh := &FormHandler{
		Config: &Config{
			Global: struct {
				Blocklist []string `toml:"blocklist" env:"BLOCKLIST" envSeparator:","`
				Port      int      `env:"PORT" envDefault:"8080"`
			}{
				Blocklist: []string{"casino", "website"},
			},
		},
	}

	tests := []struct {
		name         string
		submission   FormSubmission
		expectedPass bool
		expectedErr  error
	}{
		{
			name: "Blocklist term 'form'",
			submission: FormSubmission{
				Body: map[string]string{
					"message": "form http",
				},
				FormCfg: FormConfig{
					Blocklist: []string{"form", "http"},
					Fields: struct {
						Name     string `toml:"name"`
						Email    string `toml:"email"`
						Message  string `toml:"message"`
						Honeypot string `toml:"honeypot"`
					}{
						Message: "message",
					},
				},
			},
			expectedPass: false,
			expectedErr:  errors.New("message contains blocklist term \"form\""),
		},
		{
			name: "Blocklist term 'http'",
			submission: FormSubmission{
				Body: map[string]string{
					"message": "hello http",
				},
				FormCfg: FormConfig{
					Blocklist: []string{"test", "http"},
					Fields: struct {
						Name     string `toml:"name"`
						Email    string `toml:"email"`
						Message  string `toml:"message"`
						Honeypot string `toml:"honeypot"`
					}{
						Message: "message",
					},
				},
			},
			expectedPass: false,
			expectedErr:  errors.New("message contains blocklist term \"http\""),
		},
		{
			name: "Blocklist term 'casino'",
			submission: FormSubmission{
				Body: map[string]string{
					"message": "form http casino website",
				},
				FormCfg: FormConfig{
					Blocklist: []string{"test", "http"},
					Fields: struct {
						Name     string `toml:"name"`
						Email    string `toml:"email"`
						Message  string `toml:"message"`
						Honeypot string `toml:"honeypot"`
					}{
						Message: "message",
					},
				},
			},
			expectedPass: false,
			expectedErr:  errors.New("message contains blocklist term \"casino\""),
		},
		{
			name: "Blocklist term 'website'",
			submission: FormSubmission{
				Body: map[string]string{
					"message": "form http hello website",
				},
				FormCfg: FormConfig{
					Blocklist: []string{"test", "http"},
					Fields: struct {
						Name     string `toml:"name"`
						Email    string `toml:"email"`
						Message  string `toml:"message"`
						Honeypot string `toml:"honeypot"`
					}{
						Message: "message",
					},
				},
			},
			expectedPass: false,
			expectedErr:  errors.New("message contains blocklist term \"website\""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pass, err := c.blocklist(test.submission, *fh)
			if pass != test.expectedPass {
				t.Errorf("Expected %v, got %v", test.expectedPass, pass)
			}
			if err.Error() != test.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", test.expectedErr, err)
			}
		})
	}
}

func TestCheck_honeypot(t *testing.T) {
	c := &Check{}
	sub := FormSubmission{
		Body: map[string]string{
			"honeypot": "",
		},
		FormCfg: FormConfig{
			Fields: struct {
				Name     string "toml:\"name\""
				Email    string "toml:\"email\""
				Message  string "toml:\"message\""
				Honeypot string "toml:\"honeypot\""
			}{
				Honeypot: "honeypot",
			},
		},
	}
	expected := true
	pass, err := c.honeypot(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sub.Body["honeypot"] = "not empty"
	expected = false
	expectedErr := errors.New("honeypot field is not empty")
	pass, err = c.honeypot(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestCheck_turnstile(t *testing.T) {
	c := &Check{}
	sub := FormSubmission{
		FormCfg: FormConfig{
			TurnstileKey: "secret",
		},
		Body: map[string]string{
			"cf-turnstile-response": "token",
		},
	}
	expected := true
	pass, err := c.turnstile(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sub.FormCfg.TurnstileKey = ""
	expected = true
	pass, err = c.turnstile(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sub.FormCfg.TurnstileKey = "secret"
	sub.Body["cf-turnstile-response"] = "invalid"
	expected = false
	expectedErr := errors.New("invalid turnstile token")
	pass, err = c.turnstile(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestFormHandler_checkSpam(t *testing.T) {
	fh := &FormHandler{
		Config: &Config{
			Global: struct {
				Blocklist []string `toml:"blocklist" env:"BLOCKLIST" envSeparator:","`
				Port      int      `env:"PORT" envDefault:"8080"`
			}{
				Blocklist: []string{"global", "block"},
			},
		},
	}
	sub := FormSubmission{
		Body: map[string]string{
			"message":               "This is a test message",
			"cf-turnstile-response": "token",
			"honeypot":              "",
		},
		FormCfg: FormConfig{
			Blocklist: []string{"spam", "block"},
			Fields: struct {
				Name     string `toml:"name"`
				Email    string `toml:"email"`
				Message  string `toml:"message"`
				Honeypot string `toml:"honeypot"`
			}{
				Message:  "message",
				Honeypot: "honeypot",
			},
			TurnstileKey: "secret",
		},
	}
	expected := true
	pass := fh.checkSpam(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}

	sub.Body["honeypot"] = "not empty"
	expected = false
	pass = fh.checkSpam(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}

	sub.Body["honeypot"] = ""
	sub.Body["message"] = "This is a spam message"
	expected = false
	pass = fh.checkSpam(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}

	sub.Body["message"] = "This is a test message"
	sub.Body["cf-turnstile-response"] = "invalid"
	expected = false
	pass = fh.checkSpam(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
}
