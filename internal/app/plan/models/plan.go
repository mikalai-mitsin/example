package models

import (
	"time"

	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/uuid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Plan struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Repeat      uint64    `json:"repeat"`
	EquipmentID string    `json:"equipment_id"`
}

func (m *Plan) Validate() error {
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

type PlanFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *PlanFilter) Validate() error {
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

type PlanCreate struct {
	Name        string `json:"name"`
	Repeat      uint64 `json:"repeat"`
	EquipmentID string `json:"equipment_id"`
}

func (m *PlanCreate) Validate() error {
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

type PlanUpdate struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Repeat      *uint64   `json:"repeat"`
	EquipmentID *string   `json:"equipment_id"`
}

func (m *PlanUpdate) Validate() error {
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
