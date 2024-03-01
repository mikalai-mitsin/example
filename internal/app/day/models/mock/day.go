package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/app/day/models"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	"github.com/jaswdr/faker"
)

func NewDay(t *testing.T) *models.Day {
	t.Helper()
	return &models.Day{
		ID:          uuid.NewUUID(),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().Int(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
}
func NewDayFilter(t *testing.T) *models.DayFilter {
	t.Helper()
	return &models.DayFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewDayCreate(t *testing.T) *models.DayCreate {
	t.Helper()
	return &models.DayCreate{
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().Int(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
}
func NewDayUpdate(t *testing.T) *models.DayUpdate {
	t.Helper()
	return &models.DayUpdate{
		ID:          uuid.NewUUID(),
		Name:        pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Repeat:      pointer.Pointer(faker.New().Int()),
		EquipmentID: pointer.Pointer(faker.New().Lorem().Sentence(15)),
	}
}
