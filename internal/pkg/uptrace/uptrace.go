package uptrace

import (
	"context"

	"github.com/mikalai-mitsin/example"
	"github.com/uptrace/uptrace-go/uptrace"
)

type Config struct {
	URL         string `env:"OTEL_URL"         toml:"url"`
	Enabled     bool   `env:"OTEL_ENABLED"     toml:"enabled"`
	Environment string `env:"OTEL_ENVIRONMENT" toml:"environment"`
}
type Provider struct {
	config *Config
}

func NewProvider(config *Config) *Provider {
	return &Provider{config: config}
}
func (p *Provider) Stop(ctx context.Context) error {
	return uptrace.Shutdown(ctx)
}
func (p *Provider) Start(_ context.Context) error {
	if p.config.Enabled {
		uptrace.ConfigureOpentelemetry(
			uptrace.WithDSN(p.config.URL),
			uptrace.WithServiceName(example.Name),
			uptrace.WithServiceVersion(example.Version),
			uptrace.WithDeploymentEnvironment(p.config.Environment),
		)
	}
	return nil
}
