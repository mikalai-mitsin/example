package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/arch/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type ArchInterceptor struct {
	archUseCase ArchUseCase
	logger      log.Logger
	authUseCase AuthUseCase
}

func NewArchInterceptor(
	archUseCase ArchUseCase,
	logger log.Logger,
	authUseCase AuthUseCase,
) *ArchInterceptor {
	return &ArchInterceptor{archUseCase: archUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *ArchInterceptor) Create(
	ctx context.Context,
	create *models.ArchCreate,
) (*models.Arch, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDArchCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDArchCreate, create); err != nil {
		return nil, err
	}
	arch, err := i.archUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return arch, nil
}

func (i *ArchInterceptor) List(
	ctx context.Context,
	filter *models.ArchFilter,
) ([]*models.Arch, uint64, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDArchList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDArchList, filter); err != nil {
		return nil, 0, err
	}
	items, count, err := i.archUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *ArchInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Arch, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDArchDetail); err != nil {
		return nil, err
	}
	arch, err := i.archUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDArchDetail, arch); err != nil {
		return nil, err
	}
	return arch, nil
}

func (i *ArchInterceptor) Update(
	ctx context.Context,
	update *models.ArchUpdate,
) (*models.Arch, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDArchUpdate); err != nil {
		return nil, err
	}
	arch, err := i.archUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDArchUpdate, arch); err != nil {
		return nil, err
	}
	updated, err := i.archUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *ArchInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDArchDelete); err != nil {
		return err
	}
	arch, err := i.archUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDArchDelete, arch); err != nil {
		return err
	}
	if err := i.archUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
