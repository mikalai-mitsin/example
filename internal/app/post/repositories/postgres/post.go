package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/post/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostRepository struct {
	database *sqlx.DB
	logger   logger
}

func NewPostRepository(database *sqlx.DB, logger logger) *PostRepository {
	return &PostRepository{database: database, logger: logger}
}

type PostDTO struct {
	ID         string    `db:"id,omitempty"`
	UpdatedAt  time.Time `db:"updated_at,omitempty"`
	CreatedAt  time.Time `db:"created_at,omitempty"`
	Title      string    `db:"title"`
	Order      int64     `db:"order"`
	IsOptional bool      `db:"is_optional"`
}
type PostListDTO []*PostDTO

func (list PostListDTO) ToEntities() []*entities.Post {
	items := make([]*entities.Post, len(list))
	for i := range list {
		items[i] = list[i].ToModel()
	}
	return items
}
func NewPostDTOFromModel(model *entities.Post) *PostDTO {
	dto := &PostDTO{
		ID:         string(model.ID),
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		Title:      model.Title,
		Order:      model.Order,
		IsOptional: model.IsOptional,
	}
	return dto
}
func (dto *PostDTO) ToModel() *entities.Post {
	model := &entities.Post{
		ID:         uuid.UUID(dto.ID),
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
		Title:      dto.Title,
		Order:      dto.Order,
		IsOptional: dto.IsOptional,
	}
	return model
}
func (r *PostRepository) Create(ctx context.Context, model *entities.Post) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPostDTOFromModel(model)
	q := sq.Insert("public.posts").
		Columns("created_at", "updated_at", "title", "order", "is_optional").
		Values(dto.CreatedAt, dto.UpdatedAt, dto.Title, dto.Order, dto.IsOptional).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	model.ID = uuid.UUID(dto.ID)
	return nil
}
func (r *PostRepository) Get(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &PostDTO{}
	q := sq.Select("posts.id", "posts.created_at", "posts.updated_at", "posts.title", "posts.order", "posts.is_optional").
		From("public.posts").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("post_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *PostRepository) List(
	ctx context.Context,
	filter *entities.PostFilter,
) ([]*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto PostListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("posts.id", "posts.created_at", "posts.updated_at", "posts.title", "posts.order", "posts.is_optional").
		From("public.posts").
		Limit(pageSize)
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.ToEntities(), nil
}
func (r *PostRepository) Count(ctx context.Context, filter *entities.PostFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.posts")
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result := r.database.QueryRowxContext(ctx, query, args...)
	if err := result.Err(); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	var count uint64
	if err := result.Scan(&count); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
func (r *PostRepository) Update(ctx context.Context, model *entities.Post) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPostDTOFromModel(model)
	q := sq.Update("public.posts").Where(sq.Eq{"id": model.ID})
	{
		q = q.Set("posts.created_at", dto.CreatedAt)
		q = q.Set("posts.updated_at", dto.UpdatedAt)
		q = q.Set("posts.title", dto.Title)
		q = q.Set("posts.order", dto.Order)
		q = q.Set("posts.is_optional", dto.IsOptional)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("post_id", fmt.Sprint(model.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("post_id", fmt.Sprint(model.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("post_id", fmt.Sprint(model.ID))
		return e
	}
	return nil
}
func (r *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.posts").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("post_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("post_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("post_id", fmt.Sprint(id))
		return e
	}
	return nil
}
