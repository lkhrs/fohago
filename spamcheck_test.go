package main

import (
	"errors"
	"testing"

	"github.com/joho/godotenv"
)

func TestCheck_blocklist(t *testing.T) {
	c := &Check{}
	fh := &FormHandler{
		Config: &Config{
			Global: struct {
				Blocklist []string `env:"BLOCKLIST" envSeparator:","`
				Port      int      `env:"PORT" envDefault:"8080"`
				BaseUrl   string
			}{
				Blocklist: []string{"casino", "website"},
			},
		},
	}

	tests := []struct {
		expectedErr  error
		name         string
		submission   FormSubmission
		expectedPass bool
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
						Name     string
						Email    string
						Message  string
						Honeypot string
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
						Name     string
						Email    string
						Message  string
						Honeypot string
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
						Name     string
						Email    string
						Message  string
						Honeypot string
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
						Name     string
						Email    string
						Message  string
						Honeypot string
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
				Name     string
				Email    string
				Message  string
				Honeypot string
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

type testKey struct {
	Pass string
	Fail string
}

type testKeys struct {
	Secret testKey
	Token  string
}

var keys = testKeys{
	Secret: testKey{
		Pass: "1x0000000000000000000000000000000AA",
		Fail: "2x0000000000000000000000000000000AA",
	},
	Token: "XXXX.DUMMY.TOKEN.XXXX",
}

func TestCheck_turnstile(t *testing.T) {
	c := &Check{}
	// test pass key (true)
	sub := FormSubmission{
		FormCfg: FormConfig{
			TurnstileKey: keys.Secret.Pass,
		},
		Body: map[string]string{
			"cf-turnstile-response": keys.Token,
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

	// test blank key (true)
	sub.FormCfg.TurnstileKey = ""
	expected = true
	pass, err = c.turnstile(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// test fail key (false)
	sub.FormCfg.TurnstileKey = keys.Secret.Fail
	sub.Body["cf-turnstile-response"] = keys.Token
	expected = false
	expectedErr := errors.New("invalid-input-response")
	pass, err = c.turnstile(sub)
	if pass != expected {
		t.Errorf("Expected %v, got %v", expected, pass)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestFormHandler_checkSpam(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("Unable to load environment variables from .env")
	}
	fh := &FormHandler{
		Config: &Config{
			Global: struct {
				Blocklist []string `env:"BLOCKLIST" envSeparator:","`
				Port      int      `env:"PORT" envDefault:"8080"`
				BaseUrl   string
			}{
				Blocklist: []string{"global", "block"},
				Port:      8080,
				BaseUrl:   "http://localhost:8080",
			},
		},
	}
	sub := FormSubmission{
		Body: map[string]string{
			"message":               "This is a test message",
			"cf-turnstile-response": keys.Token,
			"honeypot":              "",
		},
		FormCfg: FormConfig{
			Blocklist: []string{"spam", "block"},
			Fields: struct {
				Name     string
				Email    string
				Message  string
				Honeypot string
			}{
				Message:  "message",
				Honeypot: "honeypot",
			},
			TurnstileKey: keys.Secret.Pass,
		},
		UserAgent: "Mozilla/5.0",
		UserIP:    "8.8.8.8",
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
		t.Errorf("Honeypot not empty: Expected %v, got %v", expected, pass)
	}

	sub.Body["honeypot"] = ""
	sub.Body["message"] = "This is a spam message"
	expected = false
	pass = fh.checkSpam(sub)
	if pass != expected {
		t.Errorf("Honeypot empty: Expected %v, got %v", expected, pass)
	}

	sub.Body["message"] = "This is a test message"
	sub.FormCfg.TurnstileKey = keys.Secret.Fail
	expected = false
	pass = fh.checkSpam(sub)
	if pass != expected {
		t.Errorf("Turnstile: Expected %v, got %v", expected, pass)
	}
}
