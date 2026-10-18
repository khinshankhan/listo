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

	// EnvLog is "json" for JSON logs, anything else for console.
	EnvLog = "LISTO_LOG"

	// EnvLogSource adds the calling file and line to each log line.
	EnvLogSource = "LISTO_LOG_SOURCE"
)

// DefaultAddr is where the server listens without EnvAddr.
const DefaultAddr = "127.0.0.1:8080"

// Addr returns the configured listen address.
func Addr() string {
	if v := strings.TrimSpace(os.Getenv(EnvAddr)); v != "" {
		return v
	}
	return DefaultAddr
}

// JSONLogs reports whether logs should be JSON.
func JSONLogs() bool {
	return strings.EqualFold(os.Getenv(EnvLog), "json")
}

// LogSource reports whether to include caller information.
func LogSource() bool {
	return truthy(os.Getenv(EnvLogSource))
}

func truthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
