package fohago

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env"
)

type Config struct {
	Api struct {
		Akismet            string `env:"AKISMET_KEY"`
		CloudFlare         string `env:"CLOUDFLARE_KEY"`
		GoogleSafeBrowsing string `env:"GOOGLE_KEY"`
	}
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
	}
}

type FormConfig struct {
	Id     string
	Fields FormFields
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
		&cfg.Api,
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
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return err
	}
	return nil
}
