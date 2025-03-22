package services

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type AuthService struct {
	authRepository authRepository
	userRepository userRepository
	logger         logger
}

func NewAuthService(
	authRepository authRepository,
	userRepository userRepository,
	logger logger,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u AuthService) CreateToken(
	ctx context.Context,
	login entities.Login,
) (entities.TokenPair, error) {
	user, err := u.userRepository.GetByEmail(ctx, login.Email)
	if err != nil {
		return entities.TokenPair{}, err
	}
	if err := user.CheckPassword(login.Password); err != nil {
		return entities.TokenPair{}, err
	}
	tokenPair, err := u.authRepository.Create(ctx, user)
	if err != nil {
		return entities.TokenPair{}, err
	}
	return tokenPair, nil
}

func (u AuthService) CreateTokenByUser(
	ctx context.Context,
	user userEntities.User,
) (entities.TokenPair, error) {
	tokenPair, err := u.authRepository.Create(ctx, user)
	if err != nil {
		return entities.TokenPair{}, err
	}
	return tokenPair, nil
}

func (u AuthService) RefreshToken(
	ctx context.Context,
	refresh entities.Token,
) (entities.TokenPair, error) {
	pair, err := u.authRepository.RefreshToken(ctx, refresh)
	if err != nil {
		return entities.TokenPair{}, err
	}
	return pair, nil
}
func (u AuthService) ValidateToken(ctx context.Context, access entities.Token) error {
	if err := u.authRepository.Validate(ctx, access); err != nil {
		return err
	}
	return nil
}
func (u AuthService) Auth(ctx context.Context, access entities.Token) (userEntities.User, error) {
	userID, err := u.authRepository.GetSubject(ctx, access)
	if err != nil {
		return userEntities.User{}, err
	}
	user, err := u.userRepository.Get(ctx, uuid.UUID(userID))
	if err != nil {
		return userEntities.User{}, err
	}
	return user, nil
}
