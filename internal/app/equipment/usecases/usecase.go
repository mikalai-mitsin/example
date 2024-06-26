package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/equipment/models"
	"github.com/018bf/example/internal/pkg/clock"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type EquipmentUseCase struct {
	equipmentRepository EquipmentRepository
	clock               clock.Clock
	logger              log.Logger
	uuid                uuid.Generator
}

func NewEquipmentUseCase(
	equipmentRepository EquipmentRepository,
	clock clock.Clock,
	logger log.Logger,
	uuid uuid.Generator,
) *EquipmentUseCase {
	return &EquipmentUseCase{
		equipmentRepository: equipmentRepository,
		clock:               clock,
		logger:              logger,
		uuid:                uuid,
	}
}

func (u *EquipmentUseCase) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
) (*models.Equipment, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	equipment := &models.Equipment{
		ID:        u.uuid.NewUUID(),
		UpdatedAt: now,
		CreatedAt: now,
		Name:      create.Name,
		Repeat:    create.Repeat,
		Weight:    create.Weight,
	}
	if err := u.equipmentRepository.Create(ctx, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (u *EquipmentUseCase) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
) ([]*models.Equipment, uint64, error) {
	equipment, err := u.equipmentRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.equipmentRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return equipment, count, nil
}
func (u *EquipmentUseCase) Get(ctx context.Context, id uuid.UUID) (*models.Equipment, error) {
	equipment, err := u.equipmentRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (u *EquipmentUseCase) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
) (*models.Equipment, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	equipment, err := u.equipmentRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.Name != nil {
			equipment.Name = *update.Name
		}
		if update.Repeat != nil {
			equipment.Repeat = *update.Repeat
		}
		if update.Weight != nil {
			equipment.Weight = *update.Weight
		}
	}
	equipment.UpdatedAt = u.clock.Now().UTC()
	if err := u.equipmentRepository.Update(ctx, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}
func (u *EquipmentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.equipmentRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
