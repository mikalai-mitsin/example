package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Comment struct {
	ID        string    `json:"id" db:"id,omitempty"`
	Body      string    `json:"body" db:"body"`
	PostId    string    `json:"post_id" db:"post_id"`
	UserId    string    `json:"user_id" db:"user_id"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty"`
}

func (c *Comment) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
		validation.Field(&c.Body),
		validation.Field(&c.PostId),
		validation.Field(&c.UserId),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type CommentFilter struct {
	IDs        []string `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
}

func (c *CommentFilter) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.IDs),
		validation.Field(&c.PageSize),
		validation.Field(&c.PageNumber),
		validation.Field(&c.OrderBy),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type CommentCreate struct {
	Body   string `json:"body"`
	PostId string `json:"post_id"`
	UserId string `json:"user_id"`
}

func (c *CommentCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Body),
		validation.Field(&c.PostId),
		validation.Field(&c.UserId),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type CommentUpdate struct {
	ID     string  `json:"id"`
	Body   *string `json:"body"`
	PostId *string `json:"post_id"`
	UserId *string `json:"user_id"`
}

func (c *CommentUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.Body),
		validation.Field(&c.PostId),
		validation.Field(&c.UserId),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDCommentList   PermissionID = "comment_list"
	PermissionIDCommentDetail PermissionID = "comment_detail"
	PermissionIDCommentCreate PermissionID = "comment_create"
	PermissionIDCommentUpdate PermissionID = "comment_update"
	PermissionIDCommentDelete PermissionID = "comment_delete"
)
