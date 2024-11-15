package mock_entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/post/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewPost(t *testing.T) *entities.Post {
	t.Helper()
	return &entities.Post{
		ID:         uuid.NewUUID(),
		CreatedAt:  faker.New().Time().Time(time.Now()),
		UpdatedAt:  faker.New().Time().Time(time.Now()),
		Title:      faker.New().Lorem().Sentence(15),
		Order:      faker.New().Int64(),
		IsOptional: faker.New().Bool(),
	}
}
func NewPostFilter(t *testing.T) *entities.PostFilter {
	t.Helper()
	return &entities.PostFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewPostCreate(t *testing.T) *entities.PostCreate {
	t.Helper()
	return &entities.PostCreate{
		Title:      faker.New().Lorem().Sentence(15),
		Order:      faker.New().Int64(),
		IsOptional: faker.New().Bool(),
	}
}
func NewPostUpdate(t *testing.T) *entities.PostUpdate {
	t.Helper()
	return &entities.PostUpdate{
		ID:         uuid.NewUUID(),
		Title:      pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Order:      pointer.Pointer(faker.New().Int64()),
		IsOptional: pointer.Pointer(faker.New().Bool()),
	}
}
