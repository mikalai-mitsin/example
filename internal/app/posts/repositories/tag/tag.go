package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagRepository struct {
	database *sqlx.DB
	logger   logger
}

func NewTagRepository(database *sqlx.DB, logger logger) *TagRepository {
	return &TagRepository{database: database, logger: logger}
}

type TagDTO struct {
	ID        uuid.UUID `db:"id,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	PostId    uuid.UUID `db:"post_id"`
	Value     string    `db:"value"`
}
type TagListDTO []TagDTO

func (list TagListDTO) ToEntities() []entities.Tag {
	items := make([]entities.Tag, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewTagDTOFromEntity(entity entities.Tag) TagDTO {
	dto := TagDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		PostId:    entity.PostId,
		Value:     entity.Value,
	}
	return dto
}
func (dto TagDTO) toEntity() entities.Tag {
	entity := entities.Tag{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		PostId:    dto.PostId,
		Value:     dto.Value,
	}
	return entity
}
func (r *TagRepository) Create(ctx context.Context, entity entities.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewTagDTOFromEntity(entity)
	q := sq.Insert("public.tags").
		Columns("id", "created_at", "updated_at", "post_id", "value").
		Values(dto.ID, dto.CreatedAt, dto.UpdatedAt, dto.PostId, dto.Value)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *TagRepository) Get(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &TagDTO{}
	q := sq.Select("tags.id", "tags.created_at", "tags.updated_at", "tags.post_id", "tags.value").
		From("public.tags").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", id.String())
		return entities.Tag{}, e
	}
	return dto.toEntity(), nil
}

func (r *TagRepository) List(
	ctx context.Context,
	filter entities.TagFilter,
) ([]entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto TagListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Of(pageSize)
	}
	q := sq.Select("tags.id", "tags.created_at", "tags.updated_at", "tags.post_id", "tags.value").
		From("public.tags").
		Limit(pageSize)
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
func (r *TagRepository) Count(ctx context.Context, filter entities.TagFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.tags")
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
func (r *TagRepository) Update(ctx context.Context, entity entities.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewTagDTOFromEntity(entity)
	q := sq.Update("public.tags").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("post_id", dto.PostId)
		q = q.Set("value", dto.Value)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("tag_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *TagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.tags").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("tag_id", fmt.Sprint(id))
		return e
	}
	return nil
}
