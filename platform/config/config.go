// Package config reads listo's environment.
package config

import (
	"os"
	"strings"
)

// Environment variables listo reads.
const (
	// EnvAddr is the address the HTTP server listens on.
	EnvAddr = "LISTO_ADDR"

	// EnvDatabaseURL is the Postgres connection string. Unset, libpq's PG*
	// variables apply.
	EnvDatabaseURL = "LISTO_DATABASE_URL"

	// EnvLog is "json" for JSON logs, anything else for console.
	EnvLog = "LISTO_LOG"

	// EnvLogSource adds the calling file and line to each log line.
	EnvLogSource = "LISTO_LOG_SOURCE"

	// EnvDev marks a development server: nothing it sends is cached.
	EnvDev = "LISTO_DEV"

	// EnvStorageURL is where images live, eg file:///var/lib/listo/media.
	EnvStorageURL = "LISTO_STORAGE_URL"
)

// DefaultAddr is where the server listens without EnvAddr.
const DefaultAddr = "127.0.0.1:8080"

// DefaultStorageURL keeps images next to the dev database, relative to the
// working directory.
const DefaultStorageURL = "file://var/media"

// Addr returns the configured listen address.
func Addr() string {
	if v := strings.TrimSpace(os.Getenv(EnvAddr)); v != "" {
		return v
	}
	return DefaultAddr
}

// DatabaseURL returns the configured connection string, which may be empty.
func DatabaseURL() string {
	return strings.TrimSpace(os.Getenv(EnvDatabaseURL))
}

// StorageURL returns the configured image store.
func StorageURL() string {
	if v := strings.TrimSpace(os.Getenv(EnvStorageURL)); v != "" {
		return v
	}
	return DefaultStorageURL
}

// JSONLogs reports whether logs should be JSON.
func JSONLogs() bool {
	return strings.EqualFold(os.Getenv(EnvLog), "json")
}

// LogSource reports whether to include caller information.
func LogSource() bool {
	return truthy(os.Getenv(EnvLogSource))
}

// Dev reports whether this is a development server.
func Dev() bool {
	return truthy(os.Getenv(EnvDev))
}

func truthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
