package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeDTO struct {
	ID        uuid.UUID  `json:"id"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	PostId    uuid.UUID  `json:"post_id"`
	Value     string     `json:"value"`
	UserId    uuid.UUID  `json:"user_id"`
}

func NewLikeDTO(entity entities.Like) (LikeDTO, error) {
	dto := LikeDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
		PostId:    entity.PostId,
		Value:     entity.Value,
		UserId:    entity.UserId,
	}
	return dto, nil
}

type LikeListDTO struct {
	Items []LikeDTO `json:"items"`
	Count uint64    `json:"count"`
}

func NewLikeListDto(likes []entities.Like, count uint64) (LikeListDTO, error) {
	response := LikeListDTO{Items: make([]LikeDTO, len(likes)), Count: count}
	for i, like := range likes {
		dto, err := NewLikeDTO(like)
		if err != nil {
			return LikeListDTO{}, err
		}
		response.Items[i] = dto
	}
	return response, nil
}

type LikeFilterDTO struct {
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	IsDeleted  *bool    `json:"is_deleted"`
}

func NewLikeFilterDTO(r *http.Request) (LikeFilterDTO, error) {
	filter := LikeFilterDTO{PageSize: nil, PageNumber: nil, OrderBy: nil, IsDeleted: nil}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return LikeFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Of(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return LikeFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_number", "Invalid page_number.").
				WithCause(err)
		}
		filter.PageNumber = pointer.Of(uint64(pageNumber))
	}
	if r.URL.Query().Has("is_deleted") {
		isDeleted, err := strconv.ParseBool(r.URL.Query().Get("is_deleted"))
		if err != nil {
			return LikeFilterDTO{}, errs.NewInvalidFormError().
				WithParam("is_deleted", "Invalid page_number.").
				WithCause(err)
		}
		filter.IsDeleted = pointer.Of(isDeleted)
	}
	if r.URL.Query().Has("order_by") {
		filter.OrderBy = strings.Split(r.URL.Query().Get("order_by"), ",")
	}
	return filter, nil
}
func (dto LikeFilterDTO) toEntity() (entities.LikeFilter, error) {
	filter := entities.LikeFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		IsDeleted:  dto.IsDeleted,
		OrderBy:    []entities.LikeOrdering{},
	}
	for _, orderBy := range dto.OrderBy {
		filter.OrderBy = append(filter.OrderBy, entities.LikeOrdering(orderBy))
	}
	return filter, nil
}

type LikeUpdateDTO struct {
	ID     uuid.UUID  `json:"id"`
	PostId *uuid.UUID `json:"post_id"`
	Value  *string    `json:"value"`
	UserId *uuid.UUID `json:"user_id"`
}

func NewLikeUpdateDTO(r *http.Request) (LikeUpdateDTO, error) {
	update := LikeUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return LikeUpdateDTO{}, err
	}
	update.ID = uuid.MustParse(chi.URLParam(r, "id"))
	return update, nil
}
func (dto LikeUpdateDTO) toEntity() (entities.LikeUpdate, error) {
	update := entities.LikeUpdate{
		ID:     dto.ID,
		PostId: dto.PostId,
		Value:  dto.Value,
		UserId: dto.UserId,
	}
	return update, nil
}

type LikeCreateDTO struct {
	PostId uuid.UUID `json:"post_id"`
	Value  string    `json:"value"`
	UserId uuid.UUID `json:"user_id"`
}

func NewLikeCreateDTO(r *http.Request) (LikeCreateDTO, error) {
	create := LikeCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return LikeCreateDTO{}, err
	}
	return create, nil
}
func (dto LikeCreateDTO) toEntity() (entities.LikeCreate, error) {
	create := entities.LikeCreate{PostId: dto.PostId, Value: dto.Value, UserId: dto.UserId}
	return create, nil
}
