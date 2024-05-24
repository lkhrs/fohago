package main

import "log"

var Conf *Config

func main() {
	Conf = &Config{}
	if err := loadFromToml(Conf, "./fohago.toml"); err != nil {
		log.Fatal(err)
	}
	if err := loadFromEnv(Conf); err != nil {
		log.Fatal(err)
	}
	if err := Conf.check(); err != nil {
		log.Fatal(err)
	}
}
