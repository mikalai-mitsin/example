package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	"github.com/jaswdr/faker"
)

func NewUser(t *testing.T) *models.User {
	t.Helper()
	return &models.User{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		FirstName: faker.New().Lorem().Sentence(15),
		LastName:  faker.New().Lorem().Sentence(15),
		Password:  faker.New().Lorem().Sentence(15),
		Email:     faker.New().Internet().Email(),
		GroupID:   models.GroupIDUser,
	}
}
func NewUserFilter(t *testing.T) *models.UserFilter {
	t.Helper()
	return &models.UserFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewUserCreate(t *testing.T) *models.UserCreate {
	t.Helper()
	return &models.UserCreate{
		FirstName: faker.New().Lorem().Sentence(15),
		LastName:  faker.New().Lorem().Sentence(15),
		Password:  faker.New().Lorem().Sentence(15),
		Email:     faker.New().Internet().Email(),
		GroupID:   models.GroupIDUser,
	}
}
func NewUserUpdate(t *testing.T) *models.UserUpdate {
	t.Helper()
	return &models.UserUpdate{
		ID:        uuid.NewUUID(),
		FirstName: pointer.Pointer(faker.New().Lorem().Sentence(15)),
		LastName:  pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Password:  pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Email:     pointer.Pointer(faker.New().Internet().Email()),
		GroupID:   pointer.Pointer(models.GroupIDUser),
	}
}
