package interceptors

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	userModels "github.com/mikalai-mitsin/example/internal/app/user/models"
)

type AuthInterceptor struct {
	authUseCase AuthUseCase
	clock       Clock
	logger      Logger
}

func NewAuthInterceptor(authUseCase AuthUseCase, clock Clock, logger Logger) *AuthInterceptor {
	return &AuthInterceptor{authUseCase: authUseCase, clock: clock, logger: logger}
}

func (i *AuthInterceptor) CreateToken(
	ctx context.Context,
	login *models.Login,
) (*models.TokenPair, error) {
	pair, err := i.authUseCase.CreateToken(ctx, login)
	if err != nil {
		return nil, err
	}
	return pair, nil
}
func (i *AuthInterceptor) ValidateToken(ctx context.Context, token models.Token) error {
	if err := i.authUseCase.ValidateToken(ctx, token); err != nil {
		return err
	}
	return nil
}

func (i *AuthInterceptor) RefreshToken(
	ctx context.Context,
	refresh models.Token,
) (*models.TokenPair, error) {
	pair, err := i.authUseCase.RefreshToken(ctx, refresh)
	if err != nil {
		return nil, err
	}
	return pair, nil
}
func (i *AuthInterceptor) Auth(ctx context.Context, access models.Token) (*userModels.User, error) {
	user, err := i.authUseCase.Auth(ctx, access)
	if err != nil {
		return nil, err
	}
	return user, nil
}
