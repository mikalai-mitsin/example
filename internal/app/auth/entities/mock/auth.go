package mock_entities // nolint:stylecheck

import (
	"testing"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"

	"github.com/jaswdr/faker"
)

func NewToken(t *testing.T) entities.Token {
	t.Helper()
	return entities.Token(faker.New().Internet().Password())
}

func NewTokenPair(t *testing.T) *entities.TokenPair {
	t.Helper()
	return &entities.TokenPair{
		Access:  NewToken(t),
		Refresh: NewToken(t),
	}
}

func NewLogin(t *testing.T) *entities.Login {
	t.Helper()
	return &entities.Login{
		Email:    faker.New().Internet().Email(),
		Password: faker.New().Internet().Password(),
	}
}
