package containers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth"
	"github.com/mikalai-mitsin/example/internal/app/post"
	"github.com/mikalai-mitsin/example/internal/app/user"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
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
}, context.Background, configs.ParseConfig, clock.NewClock, uuid.NewUUIDv4Generator, postgres.NewDatabase, postgres.NewMigrateManager, grpc.NewServer, uptrace.NewProvider, auth.NewApp, post.NewApp, user.NewApp), fx.Invoke(func(lifecycle fx.Lifecycle, server *uptrace.Provider, config *configs.Config) {
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
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, app *auth.App, server *grpc.Server) {
		lifecycle.Append(fx.Hook{OnStart: func(_ context.Context) error {
			if err := app.RegisterGRPC(server); err != nil {
				return err
			}
			return nil
		}})
	}), fx.Invoke(func(lifecycle fx.Lifecycle, app *post.App, server *grpc.Server) {
		lifecycle.Append(fx.Hook{OnStart: func(_ context.Context) error {
			if err := app.RegisterGRPC(server); err != nil {
				return err
			}
			return nil
		}})
	}), fx.Invoke(func(lifecycle fx.Lifecycle, app *user.App, server *grpc.Server) {
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
