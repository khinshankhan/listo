// Package logx is listo's logger: logstox over zap.
package logx

import (
	"time"

	"github.com/khinshankhan/logstox"
	"github.com/khinshankhan/logstox/adapter"
	"github.com/khinshankhan/logstox/backend/zapx"
	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/platform/config"
)

// Logger is what every call site holds.
type Logger = adapter.Adapter[zapx.ZapField, fields.Field]

// Field is re-exported so callers need one import.
type Field = fields.Field

// Config describes a logger.
type Config struct {
	// Name tags every line.
	Name string

	// Development selects console output over JSON.
	Development bool

	// AddSource includes the calling file and line.
	AddSource bool
}

// New builds a logger.
func New(cfg Config) Logger {
	backend := zapx.Backend{
		Development: cfg.Development,
		TimeLayout:  time.RFC3339Nano,
		AddSource:   cfg.AddSource,
		// skip the adapter frame
		CallerSkip: 2,
	}

	return adapter.Adapter[zapx.ZapField, fields.Field]{
		Base:   backend.New(logstox.Options[zapx.ZapField]{Name: cfg.Name}),
		ToBase: zapx.ToZap,
	}
}

var defaultLogger = New(Config{
	Name:        "listo",
	Development: !config.JSONLogs(),
	AddSource:   config.LogSource(),
})

// Default returns the process logger.
func Default() Logger {
	return defaultLogger
}
