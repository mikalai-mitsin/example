package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagDTO struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	PostId    uuid.UUID `json:"post_id"`
	Value     string    `json:"value"`
}

func NewTagDTO(entity entities.Tag) (TagDTO, error) {
	dto := TagDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		PostId:    entity.PostId,
		Value:     entity.Value,
	}
	return dto, nil
}

type TagListDTO struct {
	Items []TagDTO `json:"items"`
	Count uint64   `json:"count"`
}

func NewTagListDTO(tags []entities.Tag, count uint64) (TagListDTO, error) {
	response := TagListDTO{Items: make([]TagDTO, len(tags)), Count: count}
	for i, tag := range tags {
		dto, err := NewTagDTO(tag)
		if err != nil {
			return TagListDTO{}, err
		}
		response.Items[i] = dto
	}
	return response, nil
}

type TagFilterDTO struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func NewTagFilterDTO(r *http.Request) (TagFilterDTO, error) {
	filter := TagFilterDTO{IDs: nil, PageSize: nil, PageNumber: nil, OrderBy: nil}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return TagFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Of(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return TagFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_number", "Invalid page_number.").
				WithCause(err)
		}
		filter.PageNumber = pointer.Of(uint64(pageNumber))
	}
	if r.URL.Query().Has("order_by") {
		filter.OrderBy = strings.Split(r.URL.Query().Get("order_by"), ",")
	}
	if r.URL.Query().Has("ids") {
		ids := strings.Split(r.URL.Query().Get("ids"), ",")
		filter.IDs = make([]uuid.UUID, len(ids))
		for i, id := range ids {
			filter.IDs[i] = uuid.MustParse(id)
		}
	}
	return filter, nil
}
func (dto TagFilterDTO) toEntity() (entities.TagFilter, error) {
	filter := entities.TagFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		OrderBy:    dto.OrderBy,
		IDs:        dto.IDs,
	}
	return filter, nil
}

type TagUpdateDTO struct {
	ID     uuid.UUID  `json:"id"`
	PostId *uuid.UUID `json:"post_id"`
	Value  *string    `json:"value"`
}

func NewTagUpdateDTO(r *http.Request) (TagUpdateDTO, error) {
	update := TagUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return TagUpdateDTO{}, err
	}
	update.ID = uuid.MustParse(chi.URLParam(r, "id"))
	return update, nil
}
func (dto TagUpdateDTO) toEntity() (entities.TagUpdate, error) {
	update := entities.TagUpdate{ID: dto.ID, PostId: dto.PostId, Value: dto.Value}
	return update, nil
}

type TagCreateDTO struct {
	PostId uuid.UUID `json:"post_id"`
	Value  string    `json:"value"`
}

func NewTagCreateDTO(r *http.Request) (TagCreateDTO, error) {
	create := TagCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return TagCreateDTO{}, err
	}
	return create, nil
}
func (dto TagCreateDTO) toEntity() (entities.TagCreate, error) {
	create := entities.TagCreate{PostId: dto.PostId, Value: dto.Value}
	return create, nil
}
