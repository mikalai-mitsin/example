package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type UserDTO struct {
	ID        uuid.UUID        `json:"id"`
	UpdatedAt time.Time        `json:"updated_at"`
	CreatedAt time.Time        `json:"created_at"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Password  string           `json:"password"`
	Email     string           `json:"email"`
	GroupID   entities.GroupID `json:"group_id"`
}

func NewUserDTO(entity entities.User) (UserDTO, error) {
	dto := UserDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Password:  entity.Password,
		Email:     entity.Email,
		GroupID:   entity.GroupID,
	}
	return dto, nil
}

type UserListDTO struct {
	Items []UserDTO `json:"items"`
	Count uint64    `json:"count"`
}

func NewUserListDTO(users []entities.User, count uint64) (UserListDTO, error) {
	response := UserListDTO{Items: make([]UserDTO, len(users)), Count: count}
	for i, user := range users {
		dto, err := NewUserDTO(user)
		if err != nil {
			return UserListDTO{}, err
		}
		response.Items[i] = dto
	}
	return response, nil
}

type UserFilterDTO struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
	Search     string      `json:"search"`
}

func NewUserFilterDTO(r *http.Request) (UserFilterDTO, error) {
	filter := UserFilterDTO{IDs: nil, PageSize: nil, PageNumber: nil, OrderBy: nil, Search: ""}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return UserFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Pointer(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return UserFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_number", "Invalid page_number.").
				WithCause(err)
		}
		filter.PageNumber = pointer.Pointer(uint64(pageNumber))
	}
	if r.URL.Query().Has("order_by") {
		filter.OrderBy = strings.Split(r.URL.Query().Get("order_by"), ",")
	}
	if r.URL.Query().Has("ids") {
		ids := strings.Split(r.URL.Query().Get("ids"), ",")
		filter.IDs = make([]uuid.UUID, len(ids), len(ids))
		for i, id := range ids {
			filter.IDs[i] = uuid.UUID(id)
		}
	}
	if r.URL.Query().Has("search") {
		filter.Search = r.URL.Query().Get("search")
	}
	return filter, nil
}
func (dto UserFilterDTO) toEntity() (entities.UserFilter, error) {
	filter := entities.UserFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		OrderBy:    dto.OrderBy,
		IDs:        dto.IDs,
		Search:     pointer.Pointer(dto.Search),
	}
	return filter, nil
}

type UserUpdateDTO struct {
	ID        uuid.UUID         `json:"id"`
	FirstName *string           `json:"first_name"`
	LastName  *string           `json:"last_name"`
	Password  *string           `json:"password"`
	Email     *string           `json:"email"`
	GroupID   *entities.GroupID `json:"group_id"`
}

func NewUserUpdateDTO(r *http.Request) (UserUpdateDTO, error) {
	update := UserUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return UserUpdateDTO{}, err
	}
	update.ID = uuid.UUID(chi.URLParam(r, "id"))
	return update, nil
}
func (dto UserUpdateDTO) toEntity() (entities.UserUpdate, error) {
	update := entities.UserUpdate{
		ID:        dto.ID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  dto.Password,
		Email:     dto.Email,
		GroupID:   dto.GroupID,
	}
	return update, nil
}

type UserCreateDTO struct {
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Password  string           `json:"password"`
	Email     string           `json:"email"`
	GroupID   entities.GroupID `json:"group_id"`
}

func NewUserCreateDTO(r *http.Request) (UserCreateDTO, error) {
	create := UserCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return UserCreateDTO{}, err
	}
	return create, nil
}
func (dto UserCreateDTO) toEntity() (entities.UserCreate, error) {
	create := entities.UserCreate{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  dto.Password,
		Email:     dto.Email,
		GroupID:   dto.GroupID,
	}
	return create, nil
}
