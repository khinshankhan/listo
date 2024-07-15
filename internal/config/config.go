package config

// MetaConfig is the automatically filled meta information about the config
type MetaConfig struct {
	CommitHash string
	BuildDate  string
}

// Config is the main configuration structure
type Config struct {
	// Meta contains meta information about the config
	// it isn't specified in the user config
	Meta *MetaConfig
}

// Load loads a custom configuration file
func Load(commitHash, buildDate string) (*Config, error) {
	var cfg Config

	cfg.Meta = &MetaConfig{
		CommitHash: commitHash,
		BuildDate:  buildDate,
	}

	return &cfg, nil
}
