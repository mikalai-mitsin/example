package mock_entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewWidget(t *testing.T) *entities.Widget {
	t.Helper()
	return &entities.Widget{
		ID:           uuid.NewUUID(),
		CreatedAt:    faker.New().Time().Time(time.Now()),
		UpdatedAt:    faker.New().Time().Time(time.Now()),
		FormScreenId: faker.New().Lorem().Sentence(15),
		Name:         faker.New().Lorem().Sentence(15),
		Ordering:     faker.New().Int64(),
		IsOptional:   faker.New().Bool(),
		UiSettings:   faker.New().Lorem().Sentence(15),
		DeletedAt:    faker.New().Time().Time(time.Now()),
	}
}
func NewWidgetFilter(t *testing.T) *entities.WidgetFilter {
	t.Helper()
	return &entities.WidgetFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewWidgetCreate(t *testing.T) *entities.WidgetCreate {
	t.Helper()
	return &entities.WidgetCreate{
		FormScreenId: faker.New().Lorem().Sentence(15),
		Name:         faker.New().Lorem().Sentence(15),
		Ordering:     faker.New().Int64(),
		IsOptional:   faker.New().Bool(),
		UiSettings:   faker.New().Lorem().Sentence(15),
		DeletedAt:    faker.New().Time().Time(time.Now()),
	}
}
func NewWidgetUpdate(t *testing.T) *entities.WidgetUpdate {
	t.Helper()
	return &entities.WidgetUpdate{
		ID:           uuid.NewUUID(),
		FormScreenId: pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Name:         pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Ordering:     pointer.Pointer(faker.New().Int64()),
		IsOptional:   pointer.Pointer(faker.New().Bool()),
		UiSettings:   pointer.Pointer(faker.New().Lorem().Sentence(15)),
		DeletedAt:    pointer.Pointer(faker.New().Time().Time(time.Now())),
	}
}
