package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type CommentRepository struct {
	database *sqlx.DB
	logger   logger
}

func NewCommentRepository(database *sqlx.DB, logger logger) *CommentRepository {
	return &CommentRepository{database: database, logger: logger}
}

type CommentDTO struct {
	ID        string    `db:"id,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	Text      string    `db:"text"`
	AuthorId  string    `db:"author_id"`
	PostId    string    `db:"post_id"`
}
type CommentListDTO []CommentDTO

func (list CommentListDTO) ToEntities() []entities.Comment {
	items := make([]entities.Comment, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewCommentDTOFromEntity(entity entities.Comment) CommentDTO {
	dto := CommentDTO{
		ID:        string(entity.ID),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Text:      entity.Text,
		AuthorId:  string(entity.AuthorId),
		PostId:    string(entity.PostId),
	}
	return dto
}
func (dto CommentDTO) toEntity() entities.Comment {
	entity := entities.Comment{
		ID:        uuid.UUID(dto.ID),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		Text:      dto.Text,
		AuthorId:  uuid.UUID(dto.AuthorId),
		PostId:    uuid.UUID(dto.PostId),
	}
	return entity
}
func (r *CommentRepository) Create(ctx context.Context, entity entities.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewCommentDTOFromEntity(entity)
	q := sq.Insert("public.comments").
		Columns("created_at", "updated_at", "text", "author_id", "post_id").
		Values(dto.CreatedAt, dto.UpdatedAt, dto.Text, dto.AuthorId, dto.PostId)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *CommentRepository) Get(ctx context.Context, id uuid.UUID) (entities.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &CommentDTO{}
	q := sq.Select("comments.id", "comments.created_at", "comments.updated_at", "comments.text", "comments.author_id", "comments.post_id").
		From("public.comments").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("comment_id", string(id))
		return entities.Comment{}, e
	}
	return dto.toEntity(), nil
}

func (r *CommentRepository) List(
	ctx context.Context,
	filter entities.CommentFilter,
) ([]entities.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto CommentListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("comments.id", "comments.created_at", "comments.updated_at", "comments.text", "comments.author_id", "comments.post_id").
		From("public.comments").
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

func (r *CommentRepository) Count(
	ctx context.Context,
	filter entities.CommentFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.comments")
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
func (r *CommentRepository) Update(ctx context.Context, entity entities.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewCommentDTOFromEntity(entity)
	q := sq.Update("public.comments").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("text", dto.Text)
		q = q.Set("author_id", dto.AuthorId)
		q = q.Set("post_id", dto.PostId)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("comment_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("comment_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("comment_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.comments").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("comment_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("comment_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("comment_id", fmt.Sprint(id))
		return e
	}
	return nil
}
