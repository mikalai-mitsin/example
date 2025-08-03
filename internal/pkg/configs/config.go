package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
)

type otel struct {
	URL         string `env:"OTEL_URL"         toml:"url"`
	Enabled     bool   `env:"OTEL_ENABLED"     toml:"enabled"`
	Environment string `env:"OTEL_ENVIRONMENT" toml:"environment"`
}

type database struct {
	URI                string `env:"DATABASE_URI"                  toml:"uri"`
	MaxOpenConnections int    `env:"DATABASE_MAX_OPEN_CONNECTIONS" toml:"max_open_connections" env-default:"50"`
	MaxIDLEConnections int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" toml:"max_idle_connections" env-default:"10"`
}

type Config struct {
	LogLevel string       `env:"LOG_LEVEL" toml:"log_level" env-default:"debug"`
	Database database     `                toml:"database"`
	Otel     otel         `                toml:"otel"`
	HTTP     *http.Config `                toml:"http"`
	GRPC     *grpc.Config `                toml:"grpc"`
}

func ParseConfig(configPath string) (*Config, error) {
	config := &Config{}
	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, config); err != nil {
			return nil, errs.NewUnexpectedBehaviorError(err.Error())
		}
	} else {
		if err := cleanenv.ReadEnv(config); err != nil {
			return nil, errs.NewUnexpectedBehaviorError(err.Error())
		}
	}
	return config, nil
}
