package models

import (
	"time"

	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/uuid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Day struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Repeat      int       `json:"repeat"`
	EquipmentID string    `json:"equipment_id"`
}

func (m *Day) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Repeat, validation.Required),
		validation.Field(&m.EquipmentID, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type DayFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *DayFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy, validation.Required),
		validation.Field(&m.IDs, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type DayCreate struct {
	Name        string `json:"name"`
	Repeat      int    `json:"repeat"`
	EquipmentID string `json:"equipment_id"`
}

func (m *DayCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Repeat, validation.Required),
		validation.Field(&m.EquipmentID, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type DayUpdate struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Repeat      *int      `json:"repeat"`
	EquipmentID *string   `json:"equipment_id"`
}

func (m *DayUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Name),
		validation.Field(&m.Repeat),
		validation.Field(&m.EquipmentID),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
