package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "8080")
	os.Setenv("SMTP_HOST", "127.0.0.1")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USER", "user")
	os.Setenv("SMTP_PASS", "pass")
	os.Setenv("BLOCKLIST", "block1,block2")

	// Load environment variables into config
	cfg := &Config{}
	err := loadFromEnv(cfg)
	if err != nil {
		t.Errorf("Error loading configuration from environment: %v", err)
	}

	// Unset environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASS")
	os.Unsetenv("BLOCKLIST")

	// Check the loaded config fields
	if cfg.Global.Port != 8080 {
		t.Errorf("Expected Global.Port to be 8080, got %d", cfg.Global.Port)
	}
	if cfg.Smtp.Host != "127.0.0.1" {
		t.Errorf("Expected Smtp.Host to be '127.0.0.1', got %s", cfg.Smtp.Host)
	}
	if cfg.Smtp.Port != 25 {
		t.Errorf("Expected Smtp.Port to be 25, got %d", cfg.Smtp.Port)
	}
	if cfg.Smtp.User != "user" {
		t.Errorf("Expected Smtp.User to be 'user', got %s", cfg.Smtp.User)
	}
	if cfg.Smtp.Password != "pass" {
		t.Errorf("Expected Smtp.Password to be 'pass', got %s", cfg.Smtp.Password)
	}
	if len(cfg.Global.Blocklist) != 2 {
		t.Errorf("Expected Global.Blocklist to have 2 items, got %d", len(cfg.Global.Blocklist))
	}
	blocklist := []string{"block1", "block2"}
	if reflect.DeepEqual(cfg.Global.Blocklist, blocklist) != true {
		t.Errorf("Expected Global.Blocklist to contain %v, got %v", blocklist, cfg.Global.Blocklist)
	}
}

var tomlContent = `
[global]
blocklist = ["casino", "webcam"]
port = 8080
[smtp]
host = "127.0.0.1"
port = 25
[forms]
[forms.form1]
blocklist = ["business", "website"]
turnstileKey = "turnstile_key1"
[forms.form1.fields]
name = "name1"
email = "email1"
message = "message1"
honeypot = "honeypot1"
[forms.form1.mail]
recipient = "recipient@example.com"
sender = "sender@example.com"
subject = "New submission"
[forms.form2]
blocklist = ["business", "website"]
turnstileKey = "turnstile_key1"
[forms.form2.fields]
name = "name1"
email = "email1"
message = "message1"
honeypot = "honeypot1"
[forms.form2.mail]
recipient = "recipient@example.com"
sender = "sender@example.com"
subject = "New submission"
`

func TestLoadFromToml(t *testing.T) {
	// Create a temporary TOML file for testing
	tmpFile, err := os.CreateTemp("", "fohago_test_*.toml")
	if err != nil {
		t.Fatalf("Failed to create temporary TOML file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test TOML content to the temporary file
	_, err = tmpFile.WriteString(tomlContent)
	if err != nil {
		t.Fatalf("Failed to write to temporary TOML file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	configFile := tmpFile.Name()
	err = loadFromToml(cfg, configFile)
	if err != nil {
		t.Errorf("Error loading configuration from TOML file: %v", err)
	}

	checkTomlConfigFields(cfg, t)
}

func checkTomlConfigFields(cfg *Config, t *testing.T) {
	if cfg.Global.Port != 8080 {
		t.Errorf("Expected Global.Port to be 8080, got %d", cfg.Global.Port)
	}
	if cfg.Smtp.Host != "127.0.0.1" {
		t.Errorf("Expected Smtp.Host to be '127.0.0.1', got %s", cfg.Smtp.Host)
	}
	if cfg.Smtp.Port != 25 {
		t.Errorf("Expected Smtp.Port to be 25, got %d", cfg.Smtp.Port)
	}
	if len(cfg.Global.Blocklist) != 2 {
		t.Errorf("Expected Global.Blocklist to have 2 items, got %d", len(cfg.Global.Blocklist))
	}
	blocklist := []string{"casino", "webcam"}
	if reflect.DeepEqual(cfg.Global.Blocklist, blocklist) != true {
		t.Errorf("Expected Global.Blocklist to contain %v, got %v", blocklist, cfg.Global.Blocklist)
	}

	// Check if the forms are loaded correctly
	for formId := range cfg.Forms {
		if len(cfg.Forms[formId].Blocklist) != 2 {
			t.Errorf("Expected Forms.form1.Blocklist to have 2 items, got %d", len(cfg.Forms[formId].Blocklist))
		}
		blocklist = []string{"business", "website"}
		if reflect.DeepEqual(cfg.Forms[formId].Blocklist, blocklist) != true {
			t.Errorf("Expected Forms.form1.Blocklist to contain %v, got %v", blocklist, cfg.Forms[formId].Blocklist)
		}
		if cfg.Forms[formId].TurnstileKey != "turnstile_key1" {
			t.Errorf("Expected Forms.form1.TurnstileKey to be 'turnstile_key1', got %s", cfg.Forms[formId].TurnstileKey)
		}
		if cfg.Forms[formId].Fields.Name != "name1" {
			t.Errorf("Expected Forms.form1.Fields.Name to be 'name1', got %s", cfg.Forms[formId].Fields.Name)
		}
		if cfg.Forms[formId].Fields.Email != "email1" {
			t.Errorf("Expected Forms.form1.Fields.Email to be 'email1', got %s", cfg.Forms[formId].Fields.Email)
		}
		if cfg.Forms[formId].Fields.Message != "message1" {
			t.Errorf("Expected Forms.form1.Fields.Message to be 'message1', got %s", cfg.Forms[formId].Fields.Message)
		}
		if cfg.Forms[formId].Fields.Honeypot != "honeypot1" {
			t.Errorf("Expected Forms.form1.Fields.Honeypot to be 'honeypot1', got %s", cfg.Forms[formId].Fields.Honeypot)
		}

		if cfg.Forms[formId].Mail.Recipient != "recipient@example.com" {
			t.Errorf("Expected Forms.form1.Mail.Recipient to be 'recipient@example.com', got %s", cfg.Forms[formId].Mail.Recipient)
		}
		if cfg.Forms[formId].Mail.Sender != "sender@example.com" {
			t.Errorf("Expected Forms.form1.Mail.Sender to be 'sender@example.com', got %s", cfg.Forms[formId].Mail.Sender)
		}
		if cfg.Forms[formId].Mail.Subject != "New submission" {
			t.Errorf("Expected Forms.form1.Mail.Subject to be 'New submission from', got %s", cfg.Forms[formId].Mail.Subject)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary TOML file for testing
	tmpFile, err := os.CreateTemp("", "fohago_test_*.toml")
	if err != nil {
		t.Fatalf("Failed to create temporary TOML file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test TOML content to the temporary file
	_, err = tmpFile.WriteString(tomlContent)
	if err != nil {
		t.Fatalf("Failed to write to temporary TOML file: %v", err)
	}
	tmpFile.Close()

	// Set environment variables for testing
	os.Setenv("SMTP_USER", "user")
	os.Setenv("SMTP_PASS", "pass")

	// Test loading the combined configuration from TOML and environment variables
	cfg, err := loadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("Error loading configuration: %v", err)
	}

	// Unset environment variables
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASS")

	// Check fields configured by environment variables
	if cfg.Smtp.User != "user" {
		t.Errorf("Expected Smtp.User to be 'user', got %s", cfg.Smtp.User)
	}
	if cfg.Smtp.Password != "pass" {
		t.Errorf("Expected Smtp.Password to be 'pass', got %s", cfg.Smtp.Password)
	}

	checkTomlConfigFields(cfg, t)
}
