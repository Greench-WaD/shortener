package config

import "flag"

type Config struct {
	RunAddr string
	BaseURL string
}

func New() *Config {
	config := &Config{
		RunAddr: "localhost:8080",
		BaseURL: "http://localhost:8080",
	}
	config.ParseFlags()
	return config
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url")
	flag.Parse()
}
