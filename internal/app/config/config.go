package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr string
	BaseURL string
}

func New() *Config {
	config := &Config{
		RunAddr: "",
		BaseURL: "",
	}
	config.Parse()
	return config
}

func (c *Config) Parse() {
	flag.StringVar(&c.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		c.RunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		c.BaseURL = envBaseURL
	}
}
