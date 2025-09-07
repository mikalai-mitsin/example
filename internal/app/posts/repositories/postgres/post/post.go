package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostRepository struct {
	readDB  database
	writeDB database
	logger  logger
}

func NewPostRepository(readDB database, writeDB database, logger logger) *PostRepository {
	return &PostRepository{readDB: readDB, writeDB: writeDB, logger: logger}
}

var orderByMap = map[entities.PostOrdering]string{
	entities.PostOrderingUpdatedAtDESC: "posts.updated_at DESC",
	entities.PostOrderingBodyASC:       "posts.body ASC",
	entities.PostOrderingBodyDESC:      "posts.body DESC",
	entities.PostOrderingIdASC:         "posts.id ASC",
	entities.PostOrderingIdDESC:        "posts.id DESC",
	entities.PostOrderingCreatedAtASC:  "posts.created_at ASC",
	entities.PostOrderingCreatedAtDESC: "posts.created_at DESC",
	entities.PostOrderingUpdatedAtASC:  "posts.updated_at ASC",
}

func encodeOrderBy(orderBy []entities.PostOrdering) []string {
	columns := make([]string, len(orderBy))
	for i, item := range orderBy {
		column, exists := orderByMap[item]
		if !exists {
			continue
		}
		columns[i] = column
	}
	return columns
}

type PostDTO struct {
	ID        uuid.UUID `db:"id,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	Body      string    `db:"body"`
}
type PostListDTO []PostDTO

func (list PostListDTO) toEntities() []entities.Post {
	items := make([]entities.Post, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewPostDTOFromEntity(entity entities.Post) PostDTO {
	dto := PostDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Body:      entity.Body,
	}
	return dto
}
func (dto PostDTO) toEntity() entities.Post {
	entity := entities.Post{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		Body:      dto.Body,
	}
	return entity
}
func (r *PostRepository) Create(ctx context.Context, tx dtx.TX, entity entities.Post) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPostDTOFromEntity(entity)
	q := sq.Insert("public.posts").
		Columns("id", "created_at", "updated_at", "body").
		Values(dto.ID, dto.CreatedAt, dto.UpdatedAt, dto.Body)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := tx.GetSQLTx().ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *PostRepository) Get(ctx context.Context, id uuid.UUID) (entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &PostDTO{}
	q := sq.Select("posts.id", "posts.created_at", "posts.updated_at", "posts.body").
		From("public.posts").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.readDB.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("post_id", id.String())
		return entities.Post{}, e
	}
	return dto.toEntity(), nil
}

func (r *PostRepository) List(
	ctx context.Context,
	filter entities.PostFilter,
) ([]entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto PostListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Of(pageSize)
	}
	q := sq.Select("posts.id", "posts.created_at", "posts.updated_at", "posts.body").
		From("public.posts").
		Limit(pageSize)
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(encodeOrderBy(filter.OrderBy)...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.readDB.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.toEntities(), nil
}
func (r *PostRepository) Count(ctx context.Context, filter entities.PostFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.posts")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	var count uint64
	if err := r.readDB.GetContext(ctx, &count, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
func (r *PostRepository) Update(ctx context.Context, tx dtx.TX, entity entities.Post) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPostDTOFromEntity(entity)
	q := sq.Update("public.posts").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("body", dto.Body)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("post_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("post_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("post_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *PostRepository) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.posts").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
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
