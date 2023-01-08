package repositories

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewJWTAuthRepository, NewPostgresUserRepository, NewPermissionRepository, NewPostRepository, NewCommentRepository),
)
