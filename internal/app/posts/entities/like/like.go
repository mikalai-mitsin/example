package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Like struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostId    uuid.UUID `json:"post_id"`
	Value     string    `json:"value"`
	UserId    uuid.UUID `json:"user_id"`
}

func (m *Like) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
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
	if err := validation.Validate(o.String(), validation.In(LikeOrderingUpdatedAtASC.String(), LikeOrderingPostIdASC.String(), LikeOrderingUserIdASC.String(), LikeOrderingUserIdDESC.String(), LikeOrderingCreatedAtDESC.String(), LikeOrderingUpdatedAtDESC.String(), LikeOrderingPostIdDESC.String(), LikeOrderingValueASC.String(), LikeOrderingValueDESC.String(), LikeOrderingIdASC.String(), LikeOrderingIdDESC.String(), LikeOrderingCreatedAtASC.String())); err != nil {
		return err
	}
	return nil
}
func (o LikeOrdering) String() string {
	return string(o)
}

const LikeOrderingUpdatedAtDESC LikeOrdering = "-updated_at"
const LikeOrderingPostIdDESC LikeOrdering = "-post_id"
const LikeOrderingUserIdASC LikeOrdering = "user_id"
const LikeOrderingIdASC LikeOrdering = "id"
const LikeOrderingPostIdASC LikeOrdering = "post_id"
const LikeOrderingValueASC LikeOrdering = "value"
const LikeOrderingValueDESC LikeOrdering = "-value"
const LikeOrderingUserIdDESC LikeOrdering = "-user_id"
const LikeOrderingIdDESC LikeOrdering = "-id"
const LikeOrderingCreatedAtASC LikeOrdering = "created_at"
const LikeOrderingCreatedAtDESC LikeOrdering = "-created_at"
const LikeOrderingUpdatedAtASC LikeOrdering = "updated_at"

type LikeFilter struct {
	PageSize   *uint64        `json:"page_size"`
	PageNumber *uint64        `json:"page_number"`
	Search     *string        `json:"search"`
	OrderBy    []LikeOrdering `json:"order_by"`
}

func (m *LikeFilter) Validate() error {
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
