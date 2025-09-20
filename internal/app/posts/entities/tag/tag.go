package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Tag struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	PostId    uuid.UUID  `json:"post_id"`
	Value     string     `json:"value"`
}

func (m *Tag) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.DeletedAt),
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
	if err := validation.Validate(o.String(), validation.In(TagOrderingIdASC.String(), TagOrderingIdDESC.String(), TagOrderingCreatedAtASC.String(), TagOrderingCreatedAtDESC.String(), TagOrderingUpdatedAtDESC.String(), TagOrderingValueASC.String(), TagOrderingValueDESC.String(), TagOrderingUpdatedAtASC.String(), TagOrderingDeletedAtASC.String(), TagOrderingDeletedAtDESC.String(), TagOrderingPostIdASC.String(), TagOrderingPostIdDESC.String())); err != nil {
		return err
	}
	return nil
}
func (o TagOrdering) String() string {
	return string(o)
}

const TagOrderingValueASC TagOrdering = "value"
const TagOrderingValueDESC TagOrdering = "-value"
const TagOrderingIdDESC TagOrdering = "-id"
const TagOrderingUpdatedAtDESC TagOrdering = "-updated_at"
const TagOrderingDeletedAtDESC TagOrdering = "-deleted_at"
const TagOrderingPostIdDESC TagOrdering = "-post_id"
const TagOrderingIdASC TagOrdering = "id"
const TagOrderingCreatedAtASC TagOrdering = "created_at"
const TagOrderingCreatedAtDESC TagOrdering = "-created_at"
const TagOrderingUpdatedAtASC TagOrdering = "updated_at"
const TagOrderingDeletedAtASC TagOrdering = "deleted_at"
const TagOrderingPostIdASC TagOrdering = "post_id"

type TagFilter struct {
	PageSize   *uint64       `json:"page_size"`
	PageNumber *uint64       `json:"page_number"`
	Search     *string       `json:"search"`
	OrderBy    []TagOrdering `json:"order_by"`
	IsDeleted  *bool         `json:"is_deleted"`
}

func (m *TagFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy),
		validation.Field(&m.IsDeleted),
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
