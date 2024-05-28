package main

import (
	"testing"
)

func TestConfigCheck(t *testing.T) {
	// Test case 1: Valid config
	// TODO: Test cases for Form struct
	config := &Config{
		Global: struct {
			Blocklist []string `env:"BLOCKLIST" envSeparator:","`
			Port      int      `env:"PORT" envDefault:"8080"`
		}{
			Blocklist: []string{"example.com", "test.com"},
			Port:      8080,
		},
		Smtp: struct {
			User     string `env:"SMTP_USER"`
			Password string `env:"SMTP_PASS"`
			Host     string `env:"SMTP_HOST" envDefault:"localhost"`
			Port     int    `env:"SMTP_PORT" envDefault:"1025"`
		}{
			User:     "smtpuser",
			Password: "smtppass",
			Host:     "smtp.example.com",
			Port:     1025,
		},
		Api: struct {
			Akismet            string `env:"AKISMET_KEY"`
			GoogleSafeBrowsing string `env:"GOOGLE_KEY"`
		}{
			Akismet:            "akismetkey",
			GoogleSafeBrowsing: "googlekey",
		},
	}

	err := config.check()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case 2: Missing PORT
	config.Global.Port = 0
	err = config.check()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Test case 3: Missing SMTP_HOST
	config.Global.Port = 8080
	config.Smtp.Host = ""
	err = config.check()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Test case 4: Missing SMTP_PORT
	config.Smtp.Host = "smtp.example.com"
	config.Smtp.Port = 0
	err = config.check()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Test case 1: Valid environment variables
	config := &Config{}
	err := loadFromEnv(config)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case 2: Invalid environment variables
	// Set invalid environment variables here and test for error
}

func TestLoadFromToml(t *testing.T) {
	// Test case 1: Valid TOML file
	config := &Config{}
	err := loadFromToml(config, "fohago.toml")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case 2: Invalid TOML file
	// Create an invalid TOML file and test for error
}
