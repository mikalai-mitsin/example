package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/app/equipment/models"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	"github.com/jaswdr/faker"
)

func NewEquipment(t *testing.T) *models.Equipment {
	t.Helper()
	return &models.Equipment{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		Name:      faker.New().Lorem().Sentence(15),
		Repeat:    faker.New().Int(),
		Weight:    faker.New().Int(),
	}
}
func NewEquipmentFilter(t *testing.T) *models.EquipmentFilter {
	t.Helper()
	return &models.EquipmentFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewEquipmentCreate(t *testing.T) *models.EquipmentCreate {
	t.Helper()
	return &models.EquipmentCreate{
		Name:   faker.New().Lorem().Sentence(15),
		Repeat: faker.New().Int(),
		Weight: faker.New().Int(),
	}
}
func NewEquipmentUpdate(t *testing.T) *models.EquipmentUpdate {
	t.Helper()
	return &models.EquipmentUpdate{
		ID:     uuid.NewUUID(),
		Name:   pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Repeat: pointer.Pointer(faker.New().Int()),
		Weight: pointer.Pointer(faker.New().Int()),
	}
}
