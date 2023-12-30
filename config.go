package minimalmux

import (
	"log/slog"
	"os"
)

type EnvKey string

const (
	EnvDev EnvKey = "dev"
	EnvStg EnvKey = "stg"
	EnvPrd EnvKey = "prd"
)

type Config struct {
	env  EnvKey
	Port string
}

func newConfig() *Config {
	env := EnvKey(os.Getenv("APP_ENV"))
	if env == "" {
		env = EnvDev
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		env:  env,
		Port: port,
	}
}

func (c *Config) LogLevel() slog.Level {
	if c.IsDev() {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func (c *Config) IsDev() bool {
	return EnvKey(c.env) == EnvDev
}

func (c *Config) IsStg() bool {
	return EnvKey(c.env) == EnvStg
}

func (c *Config) IsPrd() bool {
	return EnvKey(c.env) == EnvPrd
}
