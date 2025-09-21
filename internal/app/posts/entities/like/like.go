package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Like struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	PostId    uuid.UUID  `json:"post_id"`
	Value     string     `json:"value"`
	UserId    uuid.UUID  `json:"user_id"`
}

func (m *Like) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.DeletedAt),
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.Value, validation.Required),
		validation.Field(&m.UserId, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type LikeOrdering string

func (o LikeOrdering) Validate() error {
	if err := validation.Validate(o.String(), validation.In(LikeOrderingCreatedAtDESC.String(), LikeOrderingUpdatedAtASC.String(), LikeOrderingDeletedAtASC.String(), LikeOrderingPostIdASC.String(), LikeOrderingValueDESC.String(), LikeOrderingUserIdASC.String(), LikeOrderingIdDESC.String(), LikeOrderingUpdatedAtDESC.String(), LikeOrderingDeletedAtDESC.String(), LikeOrderingPostIdDESC.String(), LikeOrderingValueASC.String(), LikeOrderingUserIdDESC.String(), LikeOrderingIdASC.String(), LikeOrderingCreatedAtASC.String())); err != nil {
		return err
	}
	return nil
}
func (o LikeOrdering) String() string {
	return string(o)
}

const LikeOrderingIdDESC LikeOrdering = "-id"
const LikeOrderingUpdatedAtDESC LikeOrdering = "-updated_at"
const LikeOrderingDeletedAtASC LikeOrdering = "deleted_at"
const LikeOrderingValueASC LikeOrdering = "value"
const LikeOrderingValueDESC LikeOrdering = "-value"
const LikeOrderingUserIdASC LikeOrdering = "user_id"
const LikeOrderingUserIdDESC LikeOrdering = "-user_id"
const LikeOrderingCreatedAtASC LikeOrdering = "created_at"
const LikeOrderingCreatedAtDESC LikeOrdering = "-created_at"
const LikeOrderingUpdatedAtASC LikeOrdering = "updated_at"
const LikeOrderingDeletedAtDESC LikeOrdering = "-deleted_at"
const LikeOrderingPostIdASC LikeOrdering = "post_id"
const LikeOrderingPostIdDESC LikeOrdering = "-post_id"
const LikeOrderingIdASC LikeOrdering = "id"

type LikeFilter struct {
	PageSize   *uint64        `json:"page_size"`
	PageNumber *uint64        `json:"page_number"`
	Search     *string        `json:"search"`
	OrderBy    []LikeOrdering `json:"order_by"`
	IsDeleted  *bool          `json:"is_deleted"`
}

func (m *LikeFilter) Validate() error {
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

type LikeCreate struct {
	PostId uuid.UUID `json:"post_id"`
	Value  string    `json:"value"`
	UserId uuid.UUID `json:"user_id"`
}

func (m *LikeCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PostId, validation.Required, uuid.Required),
		validation.Field(&m.Value, validation.Required),
		validation.Field(&m.UserId, validation.Required, uuid.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type LikeUpdate struct {
	ID     uuid.UUID  `json:"id"`
	PostId *uuid.UUID `json:"post_id"`
	Value  *string    `json:"value"`
	UserId *uuid.UUID `json:"user_id"`
}

func (m *LikeUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.PostId),
		validation.Field(&m.Value),
		validation.Field(&m.UserId),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type LikeDelete struct {
	ID uuid.UUID `json:"id"`
}

func (m *LikeDelete) Validate() error {
	err := validation.ValidateStruct(m, validation.Field(&m.ID, validation.Required))
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
