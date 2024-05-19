package main

import "log"

var config *Config

func main() {
	config = &Config{}
	if err := loadFromToml(config, "./fohago.toml"); err != nil {
		log.Fatal(err)
	}
	if err := loadFromEnv(config); err != nil {
		log.Fatal(err)
	}
	if err := config.check(); err != nil {
		log.Fatal(err)
	}
}
