package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type WidgetUseCase struct {
	widgetService widgetService
	logger        logger
}

func NewWidgetUseCase(widgetService widgetService, logger logger) *WidgetUseCase {
	return &WidgetUseCase{widgetService: widgetService, logger: logger}
}

func (i *WidgetUseCase) Create(
	ctx context.Context,
	create *entities.WidgetCreate,
) (*entities.Widget, error) {
	widget, err := i.widgetService.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return widget, nil
}
func (i *WidgetUseCase) Get(ctx context.Context, id uuid.UUID) (*entities.Widget, error) {
	widget, err := i.widgetService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return widget, nil
}

func (i *WidgetUseCase) List(
	ctx context.Context,
	filter *entities.WidgetFilter,
) ([]*entities.Widget, uint64, error) {
	items, count, err := i.widgetService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (i *WidgetUseCase) Update(
	ctx context.Context,
	update *entities.WidgetUpdate,
) (*entities.Widget, error) {
	updated, err := i.widgetService.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *WidgetUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.widgetService.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
