package config

import (
	"flag"
	"os"
)

// Config holds all configuration for the fuzzkit service.
type Config struct {
	Addr   string
	NoAuth bool
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		Addr:   ":8427",
		NoAuth: false,
	}
}

// Load parses flags and environment variables to build the config.
// Priority: defaults < env vars < flags.
func Load() *Config {
	c := Default()

	// Environment variables
	if addr := os.Getenv("FUZZKIT_ADDR"); addr != "" {
		c.Addr = addr
	}
	if os.Getenv("FUZZKIT_NO_AUTH") == "1" || os.Getenv("FUZZKIT_NO_AUTH") == "true" {
		c.NoAuth = true
	}

	// Flags (override env)
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.BoolVar(&c.NoAuth, "no-auth", c.NoAuth, "disable authentication (dev mode)")
	flag.Parse()

	return c
}
