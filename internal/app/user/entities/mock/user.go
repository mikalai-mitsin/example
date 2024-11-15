package mock_entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewUser(t *testing.T) *entities.User {
	t.Helper()
	return &entities.User{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		FirstName: faker.New().Lorem().Sentence(15),
		LastName:  faker.New().Lorem().Sentence(15),
		Password:  faker.New().Lorem().Sentence(15),
		Email:     faker.New().Internet().Email(),
		GroupID:   entities.GroupIDUser,
	}
}
func NewUserFilter(t *testing.T) *entities.UserFilter {
	t.Helper()
	return &entities.UserFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewUserCreate(t *testing.T) *entities.UserCreate {
	t.Helper()
	return &entities.UserCreate{
		FirstName: faker.New().Lorem().Sentence(15),
		LastName:  faker.New().Lorem().Sentence(15),
		Password:  faker.New().Lorem().Sentence(15),
		Email:     faker.New().Internet().Email(),
		GroupID:   entities.GroupIDUser,
	}
}
func NewUserUpdate(t *testing.T) *entities.UserUpdate {
	t.Helper()
	return &entities.UserUpdate{
		ID:        uuid.NewUUID(),
		FirstName: pointer.Pointer(faker.New().Lorem().Sentence(15)),
		LastName:  pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Password:  pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Email:     pointer.Pointer(faker.New().Internet().Email()),
		GroupID:   pointer.Pointer(entities.GroupIDUser),
	}
}
