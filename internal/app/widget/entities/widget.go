package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Widget struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	FormScreenId string    `json:"form_screen_id"`
	Name         string    `json:"name"`
	Ordering     int64     `json:"ordering"`
	IsOptional   bool      `json:"is_optional"`
	UiSettings   string    `json:"ui_settings"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func (m *Widget) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.FormScreenId, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Ordering, validation.Required),
		validation.Field(&m.IsOptional, validation.Required),
		validation.Field(&m.UiSettings, validation.Required),
		validation.Field(&m.DeletedAt, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type WidgetFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *WidgetFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy, validation.Required),
		validation.Field(&m.IDs, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type WidgetCreate struct {
	FormScreenId string    `json:"form_screen_id"`
	Name         string    `json:"name"`
	Ordering     int64     `json:"ordering"`
	IsOptional   bool      `json:"is_optional"`
	UiSettings   string    `json:"ui_settings"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func (m *WidgetCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.FormScreenId, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Ordering, validation.Required),
		validation.Field(&m.IsOptional, validation.Required),
		validation.Field(&m.UiSettings, validation.Required),
		validation.Field(&m.DeletedAt, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type WidgetUpdate struct {
	ID           uuid.UUID  `json:"id"`
	FormScreenId *string    `json:"form_screen_id"`
	Name         *string    `json:"name"`
	Ordering     *int64     `json:"ordering"`
	IsOptional   *bool      `json:"is_optional"`
	UiSettings   *string    `json:"ui_settings"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

func (m *WidgetUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.FormScreenId),
		validation.Field(&m.Name),
		validation.Field(&m.Ordering),
		validation.Field(&m.IsOptional),
		validation.Field(&m.UiSettings),
		validation.Field(&m.DeletedAt),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
