package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/plan/models"
	"github.com/018bf/example/internal/pkg/clock"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type PlanUseCase struct {
	planRepository PlanRepository
	clock          clock.Clock
	logger         log.Logger
	uuid           uuid.Generator
}

func NewPlanUseCase(
	planRepository PlanRepository,
	clock clock.Clock,
	logger log.Logger,
	uuid uuid.Generator,
) *PlanUseCase {
	return &PlanUseCase{planRepository: planRepository, clock: clock, logger: logger, uuid: uuid}
}
func (u *PlanUseCase) Create(ctx context.Context, create *models.PlanCreate) (*models.Plan, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	plan := &models.Plan{
		ID:          u.uuid.NewUUID(),
		UpdatedAt:   now,
		CreatedAt:   now,
		Name:        create.Name,
		Repeat:      create.Repeat,
		EquipmentID: create.EquipmentID,
	}
	if err := u.planRepository.Create(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (u *PlanUseCase) List(
	ctx context.Context,
	filter *models.PlanFilter,
) ([]*models.Plan, uint64, error) {
	plan, err := u.planRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.planRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return plan, count, nil
}
func (u *PlanUseCase) Get(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	plan, err := u.planRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return plan, nil
}
func (u *PlanUseCase) Update(ctx context.Context, update *models.PlanUpdate) (*models.Plan, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	plan, err := u.planRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.Name != nil {
			plan.Name = *update.Name
		}
		if update.Repeat != nil {
			plan.Repeat = *update.Repeat
		}
		if update.EquipmentID != nil {
			plan.EquipmentID = *update.EquipmentID
		}
	}
	plan.UpdatedAt = u.clock.Now().UTC()
	if err := u.planRepository.Update(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}
func (u *PlanUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.planRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
