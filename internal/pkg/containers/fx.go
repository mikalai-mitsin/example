package containers

import (
	"context"

	articles "github.com/mikalai-mitsin/example/internal/app/articles"
	posts "github.com/mikalai-mitsin/example/internal/app/posts"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/uptrace"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FXModule = fx.Options(fx.WithLogger(func(logger *log.Log) fxevent.Logger {
	return logger
}), fx.Provide(func(config *configs.Config) (*log.Log, error) {
	return log.NewLog(config.LogLevel)
}, context.Background, configs.ParseConfig, clock.NewClock, uuid.NewUUIDv7Generator, postgres.NewDatabase, postgres.NewMigrateManager, uptrace.NewProvider, posts.NewApp, articles.NewApp), fx.Invoke(func(lifecycle fx.Lifecycle, server *uptrace.Provider, config *configs.Config) {
	lifecycle.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}))

func NewMigrateContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger *log.Log, manager *postgres.MigrateManager, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := manager.Up(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
				}
			}()
			_ = shutdowner.Shutdown()
			return nil
		}})
	}))
	return app
}
func NewGRPCContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}, func(config *configs.Config) *grpc.Config {
		return config.GRPC
	}, grpc.NewServer), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, app *posts.App, server *grpc.Server) {
		lifecycle.Append(fx.Hook{OnStart: func(_ context.Context) error {
			if err := app.RegisterGRPC(server); err != nil {
				return err
			}
			return nil
		}})
	}), fx.Invoke(func(lifecycle fx.Lifecycle, app *articles.App, server *grpc.Server) {
		lifecycle.Append(fx.Hook{OnStart: func(_ context.Context) error {
			if err := app.RegisterGRPC(server); err != nil {
				return err
			}
			return nil
		}})
	}), fx.Invoke(func(lifecycle fx.Lifecycle, logger *log.Log, server *grpc.Server, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}, OnStop: server.Stop})
	}))
	return app
}
func NewHTTPContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}, func(config *configs.Config) *http.Config {
		return config.HTTP
	}, http.NewServer), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, app *posts.App, server *http.Server) {
		lifecycle.Append(fx.Hook{OnStart: func(_ context.Context) error {
			if err := app.RegisterHTTP(server); err != nil {
				return err
			}
			return nil
		}})
	}), fx.Invoke(func(lifecycle fx.Lifecycle, app *articles.App, server *http.Server) {
		lifecycle.Append(fx.Hook{OnStart: func(_ context.Context) error {
			if err := app.RegisterHTTP(server); err != nil {
				return err
			}
			return nil
		}})
	}), fx.Invoke(func(lifecycle fx.Lifecycle, logger *log.Log, server *http.Server, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}, OnStop: server.Stop})
	}))
	return app
}
