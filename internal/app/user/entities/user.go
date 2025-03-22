package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	GroupID   GroupID   `json:"group_id"`
}

func (m *User) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.FirstName, validation.Required),
		validation.Field(&m.LastName, validation.Required),
		validation.Field(&m.Password, validation.Required),
		validation.Field(&m.Email, validation.Required, is.EmailFormat),
		validation.Field(&m.GroupID, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
func (m *User) SetPassword(password string) {
	fromPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	m.Password = string(fromPassword)
}
func (m *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)); err != nil {
		return errs.NewInvalidParameter("email or password")
	}
	return nil
}

type UserFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *UserFilter) Validate() error {
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

type UserCreate struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Password  string  `json:"password"`
	Email     string  `json:"email"`
	GroupID   GroupID `json:"group_id"`
}

func (m *UserCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.FirstName, validation.Required),
		validation.Field(&m.LastName, validation.Required),
		validation.Field(&m.Password, validation.Required),
		validation.Field(&m.Email, validation.Required, is.EmailFormat),
		validation.Field(&m.GroupID, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type UserUpdate struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Password  *string   `json:"password"`
	Email     *string   `json:"email"`
	GroupID   *GroupID  `json:"group_id"`
}

func (m *UserUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.FirstName),
		validation.Field(&m.LastName),
		validation.Field(&m.Password),
		validation.Field(&m.Email, is.EmailFormat),
		validation.Field(&m.GroupID),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
