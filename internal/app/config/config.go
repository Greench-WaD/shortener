package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr  string
	BaseURL  string
	LogLevel string
}

func New() *Config {
	config := &Config{
		RunAddr:  "",
		BaseURL:  "",
		LogLevel: "",
	}
	config.Parse()
	return config
}

func (c *Config) Parse() {
	flag.StringVar(&c.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&c.LogLevel, "l", "Info", "log level")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		c.RunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		c.BaseURL = envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.LogLevel = envLogLevel
	}
}
