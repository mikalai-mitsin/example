package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Tag struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostId    uuid.UUID `json:"post_id"`
	Value     string    `json:"value"`
}

func (m *Tag) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.Value, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type TagOrdering string

func (o TagOrdering) Validate() error {
	if err := validation.Validate(o.String(), validation.In(TagOrderingCreatedAtASC.String(), TagOrderingUpdatedAtASC.String(), TagOrderingUpdatedAtDESC.String(), TagOrderingPostIdASC.String(), TagOrderingPostIdDESC.String(), TagOrderingValueASC.String(), TagOrderingValueDESC.String(), TagOrderingIdASC.String(), TagOrderingIdDESC.String(), TagOrderingCreatedAtDESC.String())); err != nil {
		return err
	}
	return nil
}
func (o TagOrdering) String() string {
	return string(o)
}

const TagOrderingIdASC TagOrdering = "id"
const TagOrderingCreatedAtASC TagOrdering = "created_at"
const TagOrderingCreatedAtDESC TagOrdering = "-created_at"
const TagOrderingUpdatedAtDESC TagOrdering = "-updated_at"
const TagOrderingPostIdASC TagOrdering = "post_id"
const TagOrderingPostIdDESC TagOrdering = "-post_id"
const TagOrderingValueASC TagOrdering = "value"
const TagOrderingIdDESC TagOrdering = "-id"
const TagOrderingUpdatedAtASC TagOrdering = "updated_at"
const TagOrderingValueDESC TagOrdering = "-value"

type TagFilter struct {
	PageSize   *uint64       `json:"page_size"`
	PageNumber *uint64       `json:"page_number"`
	Search     *string       `json:"search"`
	OrderBy    []TagOrdering `json:"order_by"`
}

func (m *TagFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type TagCreate struct {
	PostId uuid.UUID `json:"post_id"`
	Value  string    `json:"value"`
}

func (m *TagCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PostId, validation.Required, uuid.Required),
		validation.Field(&m.Value, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type TagUpdate struct {
	ID     uuid.UUID  `json:"id"`
	PostId *uuid.UUID `json:"post_id"`
	Value  *string    `json:"value"`
}

func (m *TagUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.PostId),
		validation.Field(&m.Value),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
