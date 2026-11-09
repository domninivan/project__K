package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbPortStr := getenv("GO_DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return nil, fmt.Errorf("config: invalid GO_DB_PORT %q", dbPortStr)
	}

	cfg := &Config{
		DBHost:     os.Getenv("GO_DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("GO_DB_USER"),
		DBPassword: os.Getenv("GO_DB_PASSWORD"),
		DBName:     os.Getenv("GO_DB_NAME"),
		ServerPort: getenv("SERVER_PORT", "8081"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	switch {
	case c.DBHost == "":
		return fmt.Errorf("config: GO_DB_HOST is required")
	case c.DBUser == "":
		return fmt.Errorf("config: GO_DB_USER is required")
	case c.DBName == "":
		return fmt.Errorf("config: GO_DB_NAME is required")
	}
	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort)),
		c.DBName,
	)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
