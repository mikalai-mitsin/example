package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	userModels "github.com/mikalai-mitsin/example/internal/app/user/models"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"

	"github.com/google/uuid"
)

const refreshAudience = "refresh"
const accessAudience = "access"

// Logger - base logger interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/logger.go . Logger
type Logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}

// Clock - clock interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/clock.go . Clock
type Clock interface {
	Now() time.Time
}

type AuthRepository struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	clock      Clock
	logger     Logger
}

func NewAuthRepository(
	config *configs.Config,
	clock Clock,
	logger Logger,
) *AuthRepository {
	private, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Auth.PrivateKey))
	if err != nil {
		panic(err)
	}
	public, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.Auth.PublicKey))
	if err != nil {
		panic(err)
	}
	return &AuthRepository{
		accessTTL:  time.Duration(config.Auth.AccessTTL) * time.Second,
		refreshTTL: time.Duration(config.Auth.RefreshTTL) * time.Second,
		publicKey:  public,
		privateKey: private,
		clock:      clock,
		logger:     logger,
	}
}

func (r *AuthRepository) Create(
	_ context.Context,
	user *userModels.User,
) (*models.TokenPair, error) {
	pair := r.createPair(string(user.ID))
	return pair, nil
}

func (r *AuthRepository) createPair(subject string) *models.TokenPair {
	now := r.clock.Now().UTC()
	accessClaims := jwt.RegisteredClaims{
		Audience:  []string{accessAudience},
		ExpiresAt: jwt.NewNumericDate(now.Add(r.accessTTL)),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "",
		NotBefore: jwt.NewNumericDate(now),
		Subject:   subject,
	}
	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), accessClaims)
	accessTokenString, err := accessToken.SignedString(r.privateKey)
	if err != nil {
		return nil
	}

	refreshClaims := jwt.RegisteredClaims{
		Audience:  []string{refreshAudience},
		ExpiresAt: jwt.NewNumericDate(now.Add(r.refreshTTL)),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "",
		NotBefore: jwt.NewNumericDate(now),
		Subject:   subject,
	}
	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(r.privateKey)
	if err != nil {
		return nil
	}
	return &models.TokenPair{
		Access:  models.Token(accessTokenString),
		Refresh: models.Token(refreshTokenString),
	}
}

func (r *AuthRepository) Validate(_ context.Context, token models.Token) error {
	jwtToken, err := r.validate(token)
	if err != nil {
		return err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if !claims.VerifyAudience(accessAudience, true) {
		return errs.NewBadTokenError()
	}
	return nil
}

func (r *AuthRepository) RefreshToken(
	_ context.Context,
	token models.Token,
) (*models.TokenPair, error) {
	jwtToken, err := r.validate(token)
	if err != nil {
		return nil, err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if !claims.VerifyAudience(refreshAudience, true) {
		return nil, errs.NewBadTokenError()
	}
	pair := r.createPair(fmt.Sprint(claims["sub"]))
	return pair, nil
}

func (r *AuthRepository) validate(token models.Token) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token.String(), r.keyFunc)
	if err != nil {
		e := errs.NewBadTokenError()
		return nil, e
	}
	return jwtToken, nil
}

func (r *AuthRepository) GetSubject(_ context.Context, token models.Token) (string, error) {
	jwtToken, err := jwt.Parse(token.String(), r.keyFunc)
	if err != nil {
		e := errs.NewError(errs.ErrorCodeUnauthenticated, "Invalid token.")
		return "", e
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	return fmt.Sprint(claims["sub"]), nil
}

func (r *AuthRepository) keyFunc(_ *jwt.Token) (interface{}, error) {
	return r.publicKey, nil
}
