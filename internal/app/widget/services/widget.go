package services

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type WidgetService struct {
	widgetRepository widgetRepository
	clock            clock
	logger           logger
	uuid             UUIDGenerator
}

func NewWidgetService(
	widgetRepository widgetRepository,
	clock clock,
	logger logger,
	uuid UUIDGenerator,
) *WidgetService {
	return &WidgetService{
		widgetRepository: widgetRepository,
		clock:            clock,
		logger:           logger,
		uuid:             uuid,
	}
}

func (u *WidgetService) Create(
	ctx context.Context,
	create *entities.WidgetCreate,
) (*entities.Widget, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	widget := &entities.Widget{
		ID:           u.uuid.NewUUID(),
		UpdatedAt:    now,
		CreatedAt:    now,
		FormScreenId: create.FormScreenId,
		Name:         create.Name,
		Ordering:     create.Ordering,
		IsOptional:   create.IsOptional,
		UiSettings:   create.UiSettings,
		DeletedAt:    create.DeletedAt,
	}
	if err := u.widgetRepository.Create(ctx, widget); err != nil {
		return nil, err
	}
	return widget, nil
}
func (u *WidgetService) Get(ctx context.Context, id uuid.UUID) (*entities.Widget, error) {
	widget, err := u.widgetRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return widget, nil
}

func (u *WidgetService) List(
	ctx context.Context,
	filter *entities.WidgetFilter,
) ([]*entities.Widget, uint64, error) {
	widget, err := u.widgetRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.widgetRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return widget, count, nil
}

func (u *WidgetService) Update(
	ctx context.Context,
	update *entities.WidgetUpdate,
) (*entities.Widget, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	widget, err := u.widgetRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.FormScreenId != nil {
			widget.FormScreenId = *update.FormScreenId
		}
		if update.Name != nil {
			widget.Name = *update.Name
		}
		if update.Ordering != nil {
			widget.Ordering = *update.Ordering
		}
		if update.IsOptional != nil {
			widget.IsOptional = *update.IsOptional
		}
		if update.UiSettings != nil {
			widget.UiSettings = *update.UiSettings
		}
		if update.DeletedAt != nil {
			widget.DeletedAt = *update.DeletedAt
		}
	}
	widget.UpdatedAt = u.clock.Now().UTC()
	if err := u.widgetRepository.Update(ctx, widget); err != nil {
		return nil, err
	}
	return widget, nil
}
func (u *WidgetService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.widgetRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
