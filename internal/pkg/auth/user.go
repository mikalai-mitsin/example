package auth

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/user/entities"
)

type ctxKey int

const UserKey ctxKey = iota + 1

func PutUser(ctx context.Context, user entities.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) entities.User {
	ctxUser := ctx.Value(UserKey)
	if ctxUser == nil {
		return entities.User{}
	}
	user, ok := ctxUser.(entities.User)
	if !ok {
		return entities.User{}
	}
	return user
}
