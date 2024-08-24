package handlers

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	userModels "github.com/mikalai-mitsin/example/internal/app/user/models"
	"github.com/mikalai-mitsin/example/internal/pkg/auth"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	headerAuthorize = "authorization"
	expectedScheme  = "bearer"
)

type AuthUseCase interface {
	Auth(context.Context, models.Token) (*userModels.User, error)
}
type AuthMiddleware struct {
	logger      *log.Log
	config      *configs.Config
	authUseCase AuthUseCase
}

func NewAuthMiddleware(
	authUseCase AuthUseCase,
	logger *log.Log,
	config *configs.Config,
) *AuthMiddleware {
	return &AuthMiddleware{authUseCase: authUseCase, logger: logger, config: config}
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
	var token models.Token
	token, err := m.authFromMD(ctx)
	if err != nil {
		return ctx, err
	}
	if token == "" {
		return auth.PutUser(ctx, models.Guest), nil
	}
	user, err := m.authUseCase.Auth(ctx, token)
	if err != nil {
		return ctx, err
	}
	newCtx := auth.PutUser(ctx, user)
	return newCtx, nil
}
func (m *AuthMiddleware) authFromMD(ctx context.Context) (models.Token, error) {
	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", nil
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(
			codes.Unauthenticated,
			"Request unauthenticated with "+expectedScheme,
		)
	}
	bearerToken := strings.TrimSpace(splits[1])
	if bearerToken == "" {
		return "", status.Errorf(
			codes.Unauthenticated,
			"Request unauthenticated with "+expectedScheme,
		)
	}
	return models.Token(splits[1]), nil
}
