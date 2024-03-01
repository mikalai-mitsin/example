package containers

import (
	"context"

	archGrpcHandlers "github.com/018bf/example/internal/app/arch/handlers/grpc"
	archInterceptors "github.com/018bf/example/internal/app/arch/interceptors"
	archPostgresRepositories "github.com/018bf/example/internal/app/arch/repositories/postgres"
	archUseCases "github.com/018bf/example/internal/app/arch/usecases"
	authGrpcHandlers "github.com/018bf/example/internal/app/auth/handlers/grpc"
	authInterceptors "github.com/018bf/example/internal/app/auth/interceptors"
	authRepositories "github.com/018bf/example/internal/app/auth/repositories/jwt"
	authUseCases "github.com/018bf/example/internal/app/auth/usecases"
	dayGrpcHandlers "github.com/018bf/example/internal/app/day/handlers/grpc"
	dayInterceptors "github.com/018bf/example/internal/app/day/interceptors"
	dayPostgresRepositories "github.com/018bf/example/internal/app/day/repositories/postgres"
	dayUseCases "github.com/018bf/example/internal/app/day/usecases"
	equipmentGrpcHandlers "github.com/018bf/example/internal/app/equipment/handlers/grpc"
	equipmentInterceptors "github.com/018bf/example/internal/app/equipment/interceptors"
	equipmentPostgresRepositories "github.com/018bf/example/internal/app/equipment/repositories/postgres"
	equipmentUseCases "github.com/018bf/example/internal/app/equipment/usecases"
	planGrpcHandlers "github.com/018bf/example/internal/app/plan/handlers/grpc"
	planInterceptors "github.com/018bf/example/internal/app/plan/interceptors"
	planPostgresRepositories "github.com/018bf/example/internal/app/plan/repositories/postgres"
	planUseCases "github.com/018bf/example/internal/app/plan/usecases"
	sessionGrpcHandlers "github.com/018bf/example/internal/app/session/handlers/grpc"
	sessionInterceptors "github.com/018bf/example/internal/app/session/interceptors"
	sessionPostgresRepositories "github.com/018bf/example/internal/app/session/repositories/postgres"
	sessionUseCases "github.com/018bf/example/internal/app/session/usecases"
	userGrpcHandlers "github.com/018bf/example/internal/app/user/handlers/grpc"
	userInterceptors "github.com/018bf/example/internal/app/user/interceptors"
	userPostgresRepositories "github.com/018bf/example/internal/app/user/repositories/postgres"
	userUseCases "github.com/018bf/example/internal/app/user/usecases"
	"github.com/018bf/example/internal/pkg/clock"
	"github.com/018bf/example/internal/pkg/configs"
	grpcInterface "github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"
	postgresInterface "github.com/018bf/example/internal/pkg/postgres"
	"github.com/018bf/example/internal/pkg/uptrace"
	uptraceInterface "github.com/018bf/example/internal/pkg/uptrace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FXModule = fx.Options(fx.WithLogger(func(logger log.Logger) fxevent.Logger {
	return logger
}), fx.Provide(func(config *configs.Config) (log.Logger, error) {
	return log.NewLog(config.LogLevel)
}, context.Background, configs.ParseConfig, clock.NewRealClock, postgresInterface.NewDatabase, postgresInterface.NewMigrateManager, grpcInterface.NewServer, grpcInterface.NewRequestIDMiddleware, grpcInterface.NewAuthMiddleware, authGrpcHandlers.NewAuthServiceServer, uptraceInterface.NewProvider, fx.Annotate(authInterceptors.NewAuthInterceptor, fx.As(new(authGrpcHandlers.AuthInterceptor)), fx.As(new(grpcInterface.AuthInterceptor))), fx.Annotate(authUseCases.NewAuthUseCase, fx.As(new(authInterceptors.AuthUseCase)), fx.As(new(sessionInterceptors.AuthUseCase)), fx.As(new(equipmentInterceptors.AuthUseCase)), fx.As(new(planInterceptors.AuthUseCase)), fx.As(new(dayInterceptors.AuthUseCase)), fx.As(new(archInterceptors.AuthUseCase)), fx.As(new(userInterceptors.AuthUseCase))), fx.Annotate(authRepositories.NewAuthRepository, fx.As(new(authUseCases.AuthRepository))), fx.Annotate(userPostgresRepositories.NewPermissionRepository, fx.As(new(authUseCases.PermissionRepository))), fx.Annotate(sessionUseCases.NewSessionUseCase, fx.As(new(sessionInterceptors.SessionUseCase))), fx.Annotate(sessionInterceptors.NewSessionInterceptor, fx.As(new(sessionGrpcHandlers.SessionInterceptor))), sessionGrpcHandlers.NewSessionServiceServer, fx.Annotate(sessionPostgresRepositories.NewSessionRepository, fx.As(new(sessionUseCases.SessionRepository))), fx.Annotate(equipmentUseCases.NewEquipmentUseCase, fx.As(new(equipmentInterceptors.EquipmentUseCase))), fx.Annotate(equipmentInterceptors.NewEquipmentInterceptor, fx.As(new(equipmentGrpcHandlers.EquipmentInterceptor))), equipmentGrpcHandlers.NewEquipmentServiceServer, fx.Annotate(equipmentPostgresRepositories.NewEquipmentRepository, fx.As(new(equipmentUseCases.EquipmentRepository))), fx.Annotate(planUseCases.NewPlanUseCase, fx.As(new(planInterceptors.PlanUseCase))), fx.Annotate(planInterceptors.NewPlanInterceptor, fx.As(new(planGrpcHandlers.PlanInterceptor))), planGrpcHandlers.NewPlanServiceServer, fx.Annotate(planPostgresRepositories.NewPlanRepository, fx.As(new(planUseCases.PlanRepository))), fx.Annotate(dayUseCases.NewDayUseCase, fx.As(new(dayInterceptors.DayUseCase))), fx.Annotate(dayInterceptors.NewDayInterceptor, fx.As(new(dayGrpcHandlers.DayInterceptor))), dayGrpcHandlers.NewDayServiceServer, fx.Annotate(dayPostgresRepositories.NewDayRepository, fx.As(new(dayUseCases.DayRepository))), fx.Annotate(archUseCases.NewArchUseCase, fx.As(new(archInterceptors.ArchUseCase))), fx.Annotate(archInterceptors.NewArchInterceptor, fx.As(new(archGrpcHandlers.ArchInterceptor))), archGrpcHandlers.NewArchServiceServer, fx.Annotate(archPostgresRepositories.NewArchRepository, fx.As(new(archUseCases.ArchRepository))), fx.Annotate(userUseCases.NewUserUseCase, fx.As(new(userInterceptors.UserUseCase))), fx.Annotate(userInterceptors.NewUserInterceptor, fx.As(new(userGrpcHandlers.UserInterceptor))), userGrpcHandlers.NewUserServiceServer, fx.Annotate(userPostgresRepositories.NewUserRepository, fx.As(new(userUseCases.UserRepository)), fx.As(new(authUseCases.UserRepository)))), fx.Invoke(func(lifecycle fx.Lifecycle, server *uptrace.Provider, config *configs.Config) {
	lifecycle.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}))

func NewMigrateContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger log.Logger, manager *postgresInterface.MigrateManager, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := manager.Up(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}})
	}))
	return app
}
func NewGRPCContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger log.Logger, server *grpcInterface.Server, shutdowner fx.Shutdowner) {
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
