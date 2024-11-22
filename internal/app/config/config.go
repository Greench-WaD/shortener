package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr         string
	BaseURL         string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
}

func New() *Config {
	config := &Config{
		RunAddr:         "",
		BaseURL:         "",
		LogLevel:        "",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}
	config.Parse()
	return config
}

func (c *Config) Parse() {
	flag.StringVar(&c.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&c.LogLevel, "l", "Info", "log level")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/short-url-db.json", "db file path")
	flag.StringVar(&c.DatabaseDSN, "d", "host=localhost user=shortener password=shortener dbname=shortener sslmode=disable", "db dsn")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.RunAddr = envRunAddr
	}
	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = envBaseURL
	}
	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		c.LogLevel = envLogLevel
	}
	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.FileStoragePath = envFileStoragePath
	}
	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		c.FileStoragePath = envDatabaseDSN
	}
}
