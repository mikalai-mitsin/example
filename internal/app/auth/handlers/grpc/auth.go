package handlers

import (
	"context"
	"strings"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type AuthServiceServer struct {
	examplepb.UnimplementedAuthServiceServer
	authUseCase authUseCase
	logger      logger
}

func NewAuthServiceServer(
	authUseCase authUseCase,
	logger logger,
) *AuthServiceServer {
	return &AuthServiceServer{authUseCase: authUseCase, logger: logger}
}

func (s AuthServiceServer) CreateToken(
	ctx context.Context,
	input *examplepb.CreateToken,
) (*examplepb.TokenPair, error) {
	login := entities.Login{
		Email:    strings.ToLower(input.GetEmail()),
		Password: input.GetPassword(),
	}
	tokenPair, err := s.authUseCase.CreateToken(ctx, login)
	if err != nil {
		return nil, err
	}
	return decodeTokenPair(tokenPair), nil
}

func (s AuthServiceServer) RefreshToken(
	ctx context.Context,
	input *examplepb.RefreshToken,
) (*examplepb.TokenPair, error) {
	tokenPair, err := s.authUseCase.RefreshToken(ctx, entities.Token(input.GetToken()))
	if err != nil {
		return nil, err
	}
	return decodeTokenPair(tokenPair), nil
}

func decodeTokenPair(pair entities.TokenPair) *examplepb.TokenPair {
	return &examplepb.TokenPair{
		Access:  pair.Access.String(),
		Refresh: pair.Refresh.String(),
	}
}
