package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/day/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type DayInterceptor struct {
	dayUseCase  DayUseCase
	logger      log.Logger
	authUseCase AuthUseCase
}

func NewDayInterceptor(
	dayUseCase DayUseCase,
	logger log.Logger,
	authUseCase AuthUseCase,
) *DayInterceptor {
	return &DayInterceptor{dayUseCase: dayUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *DayInterceptor) Create(
	ctx context.Context,
	create *models.DayCreate,
) (*models.Day, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDDayCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayCreate, create); err != nil {
		return nil, err
	}
	day, err := i.dayUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return day, nil
}

func (i *DayInterceptor) List(
	ctx context.Context,
	filter *models.DayFilter,
) ([]*models.Day, uint64, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDDayList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayList, filter); err != nil {
		return nil, 0, err
	}
	items, count, err := i.dayUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *DayInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Day, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDDayDetail); err != nil {
		return nil, err
	}
	day, err := i.dayUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDetail, day); err != nil {
		return nil, err
	}
	return day, nil
}

func (i *DayInterceptor) Update(
	ctx context.Context,
	update *models.DayUpdate,
) (*models.Day, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDDayUpdate); err != nil {
		return nil, err
	}
	day, err := i.dayUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayUpdate, day); err != nil {
		return nil, err
	}
	updated, err := i.dayUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *DayInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDDayDelete); err != nil {
		return err
	}
	day, err := i.dayUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDelete, day); err != nil {
		return err
	}
	if err := i.dayUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
