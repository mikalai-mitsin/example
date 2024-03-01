package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/plan/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type PlanInterceptor struct {
	planUseCase PlanUseCase
	logger      log.Logger
	authUseCase AuthUseCase
}

func NewPlanInterceptor(
	planUseCase PlanUseCase,
	logger log.Logger,
	authUseCase AuthUseCase,
) *PlanInterceptor {
	return &PlanInterceptor{planUseCase: planUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *PlanInterceptor) Create(
	ctx context.Context,
	create *models.PlanCreate,
) (*models.Plan, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDPlanCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDPlanCreate, create); err != nil {
		return nil, err
	}
	plan, err := i.planUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (i *PlanInterceptor) List(
	ctx context.Context,
	filter *models.PlanFilter,
) ([]*models.Plan, uint64, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDPlanList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDPlanList, filter); err != nil {
		return nil, 0, err
	}
	items, count, err := i.planUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *PlanInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDPlanDetail); err != nil {
		return nil, err
	}
	plan, err := i.planUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDPlanDetail, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (i *PlanInterceptor) Update(
	ctx context.Context,
	update *models.PlanUpdate,
) (*models.Plan, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDPlanUpdate); err != nil {
		return nil, err
	}
	plan, err := i.planUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDPlanUpdate, plan); err != nil {
		return nil, err
	}
	updated, err := i.planUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *PlanInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDPlanDelete); err != nil {
		return err
	}
	plan, err := i.planUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDPlanDelete, plan); err != nil {
		return err
	}
	if err := i.planUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
