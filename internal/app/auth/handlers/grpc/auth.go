package grpc

import (
	"context"
	"strings"

	"github.com/018bf/example/internal/app/auth/models"
	"github.com/018bf/example/internal/pkg/grpc"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
)

// AuthInterceptor - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_interceptor.go . AuthInterceptor
type AuthInterceptor interface {
	CreateToken(ctx context.Context, login *models.Login) (*models.TokenPair, error)
	RefreshToken(ctx context.Context, login models.Token) (*models.TokenPair, error)
}

type AuthServiceServer struct {
	examplepb.UnimplementedAuthServiceServer
	authInterceptor AuthInterceptor
}

func NewAuthServiceServer(
	authInterceptor AuthInterceptor,
) examplepb.AuthServiceServer {
	return &AuthServiceServer{authInterceptor: authInterceptor}
}

func (s AuthServiceServer) CreateToken(
	ctx context.Context,
	input *examplepb.CreateToken,
) (*examplepb.TokenPair, error) {
	login := &models.Login{
		Email:    strings.ToLower(input.GetEmail()),
		Password: input.GetPassword(),
	}
	tokenPair, err := s.authInterceptor.CreateToken(ctx, login)
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeTokenPair(tokenPair), nil
}

func (s AuthServiceServer) RefreshToken(
	ctx context.Context,
	input *examplepb.RefreshToken,
) (*examplepb.TokenPair, error) {
	tokenPair, err := s.authInterceptor.RefreshToken(ctx, models.Token(input.GetToken()))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeTokenPair(tokenPair), nil
}

func decodeTokenPair(pair *models.TokenPair) *examplepb.TokenPair {
	return &examplepb.TokenPair{
		Access:  pair.Access.String(),
		Refresh: pair.Refresh.String(),
	}
}
