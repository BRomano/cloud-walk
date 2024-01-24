package main

import (
	"fmt"
	"github.com/vrischmann/envconfig"
	"time"
)

type Settings struct {
	Server struct {
		Port            string        `envconfig:"default=8080"`
		Context         string        `envconfig:"default=log-parser"`
		ReadTimeout     time.Duration `envconfig:"default=3s"`
		WriteTimeout    time.Duration `envconfig:"default=4s"`
		ShutdownTimeout time.Duration `envconfig:"default=15s"`
		LogLevel        string        `envconfig:"default=TRACE"`

		RequestTimeout time.Duration `envconfig:"default=1.5s"`
	}
}

func (settings *Settings) initSettings() error {
	options := envconfig.Options{AllOptional: false}
	if err := envconfig.InitWithOptions(settings, options); err != nil {
		return fmt.Errorf("initializing settings: %w", err)
	}
	return nil
}
