package handlers

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/auth"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"google.golang.org/grpc"
)

const (
	headerAuthorize = "authorization"
	expectedScheme  = "bearer"
)

type AuthService interface {
	Auth(context.Context, entities.Token) (userEntities.User, error)
}
type AuthMiddleware struct {
	logger      *log.Log
	config      *configs.Config
	authService AuthService
}

func NewAuthMiddleware(
	authService AuthService,
	logger *log.Log,
	config *configs.Config,
) *AuthMiddleware {
	return &AuthMiddleware{authService: authService, logger: logger, config: config}
}

func (m *AuthMiddleware) UnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	newCtx, err := m.auth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(newCtx, req)
}
func (m *AuthMiddleware) auth(ctx context.Context) (context.Context, error) {
	var token entities.Token
	token, err := m.authFromMD(ctx)
	if err != nil {
		return ctx, err
	}
	if token == "" {
		return auth.PutUser(ctx, entities.Guest), nil
	}
	user, err := m.authService.Auth(ctx, token)
	if err != nil {
		return ctx, err
	}
	newCtx := auth.PutUser(ctx, user)
	return newCtx, nil
}
func (m *AuthMiddleware) authFromMD(ctx context.Context) (entities.Token, error) {
	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", nil
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", errs.NewUnauthenticatedError()
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", errs.NewUnauthenticatedError()
	}
	bearerToken := strings.TrimSpace(splits[1])
	if bearerToken == "" {
		return "", errs.NewUnauthenticatedError()
	}
	return entities.Token(splits[1]), nil
}
