package handlers

import (
	"context"
	"strings"

	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type AuthServiceServer struct {
	examplepb.UnimplementedAuthServiceServer
	authInterceptor AuthInterceptor
}

func NewAuthServiceServer(
	authInterceptor AuthInterceptor,
) *AuthServiceServer {
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
		return nil, err
	}
	return decodeTokenPair(tokenPair), nil
}

func (s AuthServiceServer) RefreshToken(
	ctx context.Context,
	input *examplepb.RefreshToken,
) (*examplepb.TokenPair, error) {
	tokenPair, err := s.authInterceptor.RefreshToken(ctx, models.Token(input.GetToken()))
	if err != nil {
		return nil, err
	}
	return decodeTokenPair(tokenPair), nil
}

func decodeTokenPair(pair *models.TokenPair) *examplepb.TokenPair {
	return &examplepb.TokenPair{
		Access:  pair.Access.String(),
		Refresh: pair.Refresh.String(),
	}
}
