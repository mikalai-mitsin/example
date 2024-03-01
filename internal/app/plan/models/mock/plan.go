package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/app/plan/models"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	"github.com/jaswdr/faker"
)

func NewPlan(t *testing.T) *models.Plan {
	t.Helper()
	return &models.Plan{
		ID:          uuid.NewUUID(),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().UInt64(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
}
func NewPlanFilter(t *testing.T) *models.PlanFilter {
	t.Helper()
	return &models.PlanFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewPlanCreate(t *testing.T) *models.PlanCreate {
	t.Helper()
	return &models.PlanCreate{
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().UInt64(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
}
func NewPlanUpdate(t *testing.T) *models.PlanUpdate {
	t.Helper()
	return &models.PlanUpdate{
		ID:          uuid.NewUUID(),
		Name:        pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Repeat:      pointer.Pointer(faker.New().UInt64()),
		EquipmentID: pointer.Pointer(faker.New().Lorem().Sentence(15)),
	}
}
