package containers

import (
	"context"
	postgresInterface "github.com/018bf/example/internal/interfaces/postgres"
	jwtRepositories "github.com/018bf/example/internal/repositories/jwt"
	postgresRepositories "github.com/018bf/example/internal/repositories/postgres"

	"github.com/018bf/example/pkg/log"
	"go.uber.org/fx/fxevent"

	"github.com/018bf/example/internal/interceptors"
	"github.com/018bf/example/internal/usecases"

	"github.com/018bf/example/pkg/clock"

	"github.com/018bf/example/internal/configs"
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.WithLogger(
		func(logger log.Logger) fxevent.Logger {
			return logger
		},
	),
	fx.Provide(
		context.Background,
		configs.ParseConfig,
		clock.NewRealClock,
		postgresInterface.NewDatabase,
		postgresInterface.NewMigrateManager,
		func(config *configs.Config) (log.Logger, error) {
			return log.NewLog(config.LogLevel)
		},
		usecases.NewAuthUseCase,
		interceptors.NewAuthInterceptor,
		jwtRepositories.NewAuthRepository,
		postgresRepositories.NewPermissionRepository,
		usecases.NewUserUseCase,
		interceptors.NewUserInterceptor,
		postgresRepositories.NewPostgresUserRepository, usecases.NewSessionUseCase, interceptors.NewSessionInterceptor, postgresRepositories.NewSessionRepository, usecases.NewEquipmentUseCase, interceptors.NewEquipmentInterceptor, postgresRepositories.NewEquipmentRepository, usecases.NewPlanUseCase, interceptors.NewPlanInterceptor, postgresRepositories.NewPlanRepository, usecases.NewDayUseCase, interceptors.NewDayInterceptor, postgresRepositories.NewDayRepository, usecases.NewArchUseCase, interceptors.NewArchInterceptor, postgresRepositories.NewArchRepository,
	),
)

func NewMigrate(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string { return config }),
		FXModule,
		fx.Invoke(func(lifecycle fx.Lifecycle, manager *postgresInterface.MigrateManager) {
			lifecycle.Append(fx.Hook{
				OnStart: manager.Up,
			})
		}),
	)
	return app
}
