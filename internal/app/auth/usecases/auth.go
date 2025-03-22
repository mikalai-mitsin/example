package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
)

type AuthUseCase struct {
	authService authService
	clock       clock
	logger      logger
}

func NewAuthUseCase(authService authService, clock clock, logger logger) *AuthUseCase {
	return &AuthUseCase{authService: authService, clock: clock, logger: logger}
}

func (i *AuthUseCase) CreateToken(
	ctx context.Context,
	login entities.Login,
) (entities.TokenPair, error) {
	pair, err := i.authService.CreateToken(ctx, login)
	if err != nil {
		return entities.TokenPair{}, err
	}
	return pair, nil
}
func (i *AuthUseCase) ValidateToken(ctx context.Context, token entities.Token) error {
	if err := i.authService.ValidateToken(ctx, token); err != nil {
		return err
	}
	return nil
}

func (i *AuthUseCase) RefreshToken(
	ctx context.Context,
	refresh entities.Token,
) (entities.TokenPair, error) {
	pair, err := i.authService.RefreshToken(ctx, refresh)
	if err != nil {
		return entities.TokenPair{}, err
	}
	return pair, nil
}
func (i *AuthUseCase) Auth(ctx context.Context, access entities.Token) (userEntities.User, error) {
	user, err := i.authService.Auth(ctx, access)
	if err != nil {
		return userEntities.User{}, err
	}
	return user, nil
}
