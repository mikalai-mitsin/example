package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/equipment/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type EquipmentInterceptor struct {
	equipmentUseCase EquipmentUseCase
	logger           log.Logger
	authUseCase      AuthUseCase
}

func NewEquipmentInterceptor(
	equipmentUseCase EquipmentUseCase,
	logger log.Logger,
	authUseCase AuthUseCase,
) *EquipmentInterceptor {
	return &EquipmentInterceptor{
		equipmentUseCase: equipmentUseCase,
		logger:           logger,
		authUseCase:      authUseCase,
	}
}

func (i *EquipmentInterceptor) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
) (*models.Equipment, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate, create); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
) ([]*models.Equipment, uint64, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentList, filter); err != nil {
		return nil, 0, err
	}
	items, count, err := i.equipmentUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *EquipmentInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Equipment, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
) (*models.Equipment, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate, equipment); err != nil {
		return nil, err
	}
	updated, err := i.equipmentUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *EquipmentInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete); err != nil {
		return err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete, equipment); err != nil {
		return err
	}
	if err := i.equipmentUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
