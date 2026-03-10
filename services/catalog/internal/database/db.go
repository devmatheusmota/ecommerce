package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Open() (*sql.DB, error) {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		connectionString = buildConnectionString()
	}

	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return database, nil
}

func buildConnectionString() string {
	host := getEnv("PG_HOST", "localhost")
	port := getEnv("PG_PORT", "5432")
	user := getEnv("PG_USER", "ecommerce")
	password := getEnv("PG_PASSWORD", "ecommerce_dev")
	database := getEnv("PG_DATABASE", "catalog")
	sslMode := getEnv("PG_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, database, sslMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
