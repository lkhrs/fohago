package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env"
)

type Config struct {
	Forms map[string]FormConfig
	Smtp  struct {
		User     string `env:"SMTP_USER"`
		Password string `env:"SMTP_PASS"`
		Host     string `env:"SMTP_HOST" envDefault:"localhost"`
		Port     int    `env:"SMTP_PORT" envDefault:"1025"`
	}
	Global struct {
		Blocklist []string `env:"BLOCKLIST" envSeparator:","`
		Port      int      `env:"PORT" envDefault:"8080"`
		BaseUrl   string
		LogLevel  string
	}
}

type FormBody map[string]string

type FormConfig struct {
	Id   string
	Body FormBody
	Mail struct {
		Recipient string
		Sender    string
		Subject   string
	}
	TurnstileKey string
	Fields       struct {
		Name     string
		Email    string
		Message  string
		Honeypot string
	}
	Blocklist []string
	Redirects struct {
		Success string
	}
}

// check the config for required fields
func (c *Config) check() error {
	if c.Global.Port == 0 {
		return fmt.Errorf("PORT is required")
	}
	if c.Smtp.Host == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.Smtp.Port == 0 {
		return fmt.Errorf("SMTP_PORT is required")
	}
	return nil
}

// loads the configuration from the environment
func loadFromEnv(cfg *Config) error {
	fields := []interface{}{
		&cfg.Global,
		&cfg.Smtp,
	}
	for _, field := range fields {
		if err := env.Parse(field); err != nil {
			return err
		}
	}
	return nil
}

// loads the configuration from a TOML file
func loadFromToml(cfg *Config, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(file string) *Config {
	cfg := &Config{}

	if err := loadFromEnv(cfg); err != nil {
		slog.Warn("Could not load environment variables:", slog.Any("error", err))
		slog.Info("Using defaults")
	}

	if err := loadFromToml(cfg, file); err != nil && !os.IsNotExist(err) {
		slog.Warn("Problem loading config file:", slog.Any("error", err))
	} else if os.IsNotExist(err) {
		slog.Warn("Could not load config file:", slog.Any("error", err))
	}

	if err := cfg.check(); err != nil {
		slog.Error("Missing required option(s):", slog.Any("error", err))
	}

	return cfg
}
