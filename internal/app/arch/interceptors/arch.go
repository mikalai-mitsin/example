package interceptors

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/arch/models"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArchInterceptor struct {
	archUseCase ArchUseCase
	logger      Logger
}

func NewArchInterceptor(archUseCase ArchUseCase, logger Logger) *ArchInterceptor {
	return &ArchInterceptor{archUseCase: archUseCase, logger: logger}
}

func (i *ArchInterceptor) Create(
	ctx context.Context,
	create *models.ArchCreate,
) (*models.Arch, error) {
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
	items, count, err := i.archUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *ArchInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Arch, error) {
	arch, err := i.archUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return arch, nil
}

func (i *ArchInterceptor) Update(
	ctx context.Context,
	update *models.ArchUpdate,
) (*models.Arch, error) {
	updated, err := i.archUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *ArchInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.archUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
