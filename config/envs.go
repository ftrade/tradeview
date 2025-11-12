package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Envs struct {
	FontFilePath   string `env:"FONT_FILE_PATH,notEmpty"`
	FontSize       int32  `env:"FONT_SIZE"                 envDefault:"20"`
	MarketFilePath string `env:"MARKET_FILE_PATH,notEmpty"`
	VsyncEnabled   bool   `env:"VSYNC_ENABLED"`
	LogLevel       string `env:"LOG_LEVEL,notEmpty"        envDefault:"INFO"`
}

func LoadEnvs(fileNames ...string) (Envs, error) {
	var envs Envs
	if err := godotenv.Load(fileNames...); err != nil {
		return envs, fmt.Errorf("failed to load env file: %w", err)
	}

	if err := env.Parse(&envs); err != nil {
		return envs, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}

	return envs, nil
}
