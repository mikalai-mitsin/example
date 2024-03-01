package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/app/session/models"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	"github.com/jaswdr/faker"
)

func NewSession(t *testing.T) *models.Session {
	t.Helper()
	return &models.Session{
		ID:          uuid.NewUUID(),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		Title:       faker.New().Lorem().Sentence(15),
		Description: faker.New().Lorem().Sentence(15),
	}
}
func NewSessionFilter(t *testing.T) *models.SessionFilter {
	t.Helper()
	return &models.SessionFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewSessionCreate(t *testing.T) *models.SessionCreate {
	t.Helper()
	return &models.SessionCreate{
		Title:       faker.New().Lorem().Sentence(15),
		Description: faker.New().Lorem().Sentence(15),
	}
}
func NewSessionUpdate(t *testing.T) *models.SessionUpdate {
	t.Helper()
	return &models.SessionUpdate{
		ID:          uuid.NewUUID(),
		Title:       pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Description: pointer.Pointer(faker.New().Lorem().Sentence(15)),
	}
}
