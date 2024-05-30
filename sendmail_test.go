package main

import (
	"testing"
)

func TestBuildAndSend(t *testing.T) {
	cfg := &Config{
		Smtp: struct {
			User     string `env:"SMTP_USER"`
			Password string `env:"SMTP_PASS"`
			Host     string `env:"SMTP_HOST" envDefault:"localhost"`
			Port     int    `env:"SMTP_PORT" envDefault:"1025"`
		}{
			User:     "",
			Password: "",
			Host:     "localhost",
			Port:     1025,
		},
	}

	sub := FormSubmission{
		Id: "example",
		FormCfg: FormConfig{
			Mail: struct {
				Recipient string `toml:"recipient"`
				Sender    string `toml:"sender"`
				Subject   string `toml:"subject"`
			}{
				Recipient: "recipient@example.com",
				Sender:    "sender@example.com",
				Subject:   "Test Subject",
			},
			Fields: struct {
				Name     string `toml:"name"`
				Email    string `toml:"email"`
				Message  string `toml:"message"`
				Honeypot string `toml:"honeypot"`
			}{
				Name:     "name",
				Email:    "email",
				Message:  "message",
				Honeypot: "honeypot",
			},
		},
		Body: map[string]string{
			"email": "replyto@example.com",
		},
	}

	err := buildAndSend(cfg, sub)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestBuildEmailMessage(t *testing.T) {
	formCfg := FormConfig{
		Mail: struct {
			Recipient string `toml:"recipient"`
			Sender    string `toml:"sender"`
			Subject   string `toml:"subject"`
		}{
			Recipient: "recipient@example.com",
			Sender:    "sender@example.com",
			Subject:   "Test Subject",
		},
		Fields: struct {
			Name     string `toml:"name"`
			Email    string `toml:"email"`
			Message  string `toml:"message"`
			Honeypot string `toml:"honeypot"`
		}{
			Name:     "name",
			Email:    "email",
			Message:  "message",
			Honeypot: "honeypot",
		},
	}

	sub := FormSubmission{
		Id:      "example",
		FormCfg: formCfg,
		Body: map[string]string{
			"name":    "TestName",
			"email":   "test@example.com",
			"message": "Testing the message field.",
		},
	}

	msg, err := buildEmailMessage(sub)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	headers := "From: <" + sub.FormCfg.Mail.Sender + ">\r\n" +
		"To: <" + sub.FormCfg.Mail.Recipient + ">\r\n" +
		"Subject: " + sub.FormCfg.Mail.Subject + " - " + sub.Id + "\r\n" +
		"Reply-To: <" + sub.Body[sub.FormCfg.Fields.Email] + ">\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	expectedMsg := message{
		Subject:   "Test Subject - example",
		Body:      []byte(headers + "TestName test@example.com Testing the message field."),
		Recipient: "<recipient@example.com>",
		Sender:    "<sender@example.com>",
		ReplyTo:   "test@example.com",
	}

	if string(msg.Body) != string(expectedMsg.Body) {
		t.Errorf("Expected body %q, got %q", expectedMsg.Body, msg.Body)
	}

	if msg.Recipient != expectedMsg.Recipient {
		t.Errorf("Expected recipient %q, got %q", expectedMsg.Recipient, msg.Recipient)
	}

	if msg.Sender != expectedMsg.Sender {
		t.Errorf("Expected sender %q, got %q", expectedMsg.Sender, msg.Sender)
	}

	if msg.ReplyTo != expectedMsg.ReplyTo {
		t.Errorf("Expected reply-to %q, got %q", expectedMsg.ReplyTo, msg.ReplyTo)
	}
}

func TestSendEmail(t *testing.T) {
	cfg := &Config{
		Smtp: struct {
			User     string `env:"SMTP_USER"`
			Password string `env:"SMTP_PASS"`
			Host     string `env:"SMTP_HOST" envDefault:"localhost"`
			Port     int    `env:"SMTP_PORT" envDefault:"1025"`
		}{
			User:     "",
			Password: "",
			Host:     "localhost",
			Port:     1025,
		},
	}

	msg := message{
		Subject:   "Test Subject",
		Body:      []byte("Test message"),
		Recipient: "recipient@example.com",
		Sender:    "sender@example.com",
		ReplyTo:   "replyto@example.com",
	}

	err := sendEmail(cfg, msg)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestLoadTemplate(t *testing.T) {
	id := "example"

	tmpl := loadTemplate(id)
	if tmpl == nil {
		t.Errorf("Expected template, got nil")
	}
}
