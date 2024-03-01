package models

import (
	"time"

	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/uuid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Session struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (m *Session) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *SessionFilter) Validate() error {
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

type SessionCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (m *SessionCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Description, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type SessionUpdate struct {
	ID          uuid.UUID `json:"id"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
}

func (m *SessionUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Title),
		validation.Field(&m.Description),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
