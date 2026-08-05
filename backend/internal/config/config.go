// Package config holds the runtime configuration for the API server.
package config

import "time"

// Config is the top-level application configuration, populated from flags/env
// in cmd/api/main.go and injected into the application struct.
type Config struct {
	Port int
	Env  string
	DB   DB
	JWT  JWT
}

// DB holds database connection and pool settings.
type DB struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

// JWT holds token signing settings.
type JWT struct {
	Secret string
	Issuer string
	TTL    time.Duration
}
