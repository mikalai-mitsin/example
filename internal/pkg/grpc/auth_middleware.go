package grpc

import (
	"context"
	"strings"

	"github.com/018bf/example/internal/app/auth/models"
	"github.com/018bf/example/internal/pkg/configs"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ctxKey int

const UserKey ctxKey = iota + 1
const (
	headerAuthorize = "authorization"
	expectedScheme  = "bearer"
)

func AuthFromMD(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", status.Errorf(
			codes.Unauthenticated,
			"Request unauthenticated with "+expectedScheme,
		)
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
	return splits[1], nil
}

type AuthMiddleware struct {
	logger          log.Logger
	config          *configs.Config
	authInterceptor AuthInterceptor
}

func NewAuthMiddleware(
	authInterceptor AuthInterceptor,
	logger log.Logger,
	config *configs.Config,
) *AuthMiddleware {
	return &AuthMiddleware{authInterceptor: authInterceptor, logger: logger, config: config}
}
func (m *AuthMiddleware) Auth(ctx context.Context) (context.Context, error) {
	var token models.Token
	stringToken, err := AuthFromMD(ctx)
	if err != nil {
		return context.WithValue(ctx, UserKey, models.Guest), nil
	}
	if stringToken == "" {
		return context.WithValue(ctx, UserKey, models.Guest), nil
	}
	token = models.Token(stringToken)
	user, err := m.authInterceptor.Auth(ctx, token)
	if err != nil {
		return nil, DecodeError(err)
	}
	newCtx := context.WithValue(ctx, UserKey, user)
	return newCtx, nil
}

func (m *AuthMiddleware) UnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	newCtx, err := m.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(newCtx, req)
}
