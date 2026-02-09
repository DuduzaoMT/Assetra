package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	dbprovider = "postgres"
)

type Config interface {
	Dsn() string
}

// generic config struct to hold db config values
type config struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
	dbName     string
	dsn        string
	// sslMode    string
}

func NewConfig() Config {
	// load env variables
	var cfg config
	cfg.dbUser = os.Getenv("DB_USER")
	cfg.dbPassword = os.Getenv("DB_PASSWORD")
	cfg.dbHost = os.Getenv("DB_HOST")
	cfg.dbName = os.Getenv("DB_NAME")

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(fmt.Errorf("invalid DB_PORT: %v", err))
	}
	cfg.dbPort = port

	// validate env variables
	if cfg.dbUser == "" || cfg.dbPassword == "" || cfg.dbHost == "" || cfg.dbPort == 0 || cfg.dbName == "" {
		log.Fatal("missing required database environment variables")
	}
	// get sslmode from env, default to 'require'
	sslMode := "disable"
	if os.Getenv("ENV") == "production" {
		sslMode = "require"
	}
	cfg.dsn = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		dbprovider, cfg.dbUser, cfg.dbPassword, cfg.dbHost, cfg.dbPort, cfg.dbName, sslMode)

	log.Println("Database configuration loaded successfully")
	return &cfg
}

func (c *config) Dsn() string {
	return c.dsn
}
