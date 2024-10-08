package auth

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/user/models"
)

type ctxKey int

const UserKey ctxKey = iota + 1

func PutUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) *models.User {
	ctxUser := ctx.Value(UserKey)
	if ctxUser == nil {
		return nil
	}
	user, ok := ctxUser.(*models.User)
	if !ok {
		return nil
	}
	return user
}
