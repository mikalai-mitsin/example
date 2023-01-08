package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/018bf/example/pkg/log"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewCommentRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.CommentRepository {
	return &CommentRepository{
		database: database,
		logger:   logger,
	}
}

func (r *CommentRepository) Create(
	ctx context.Context,
	comment *models.Comment,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.comments").
		Columns(
			"body",
			"post_id",
			"user_id",
			"updated_at",
			"created_at",
		).
		Values(
			comment.Body,
			comment.PostId,
			comment.UserId,
			comment.UpdatedAt,
			comment.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(comment); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *CommentRepository) Get(
	ctx context.Context,
	id string,
) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	comment := &models.Comment{}
	q := sq.Select(
		"comments.id",
		"comments.body",
		"comments.post_id",
		"comments.user_id",
		"comments.updated_at",
		"comments.created_at",
	).
		From("public.comments").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, comment, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("comment_id", id)
		return nil, e
	}
	return comment, nil
}

func (r *CommentRepository) List(
	ctx context.Context,
	filter *models.CommentFilter,
) ([]*models.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var comments []*models.Comment
	const pageSize = 10
	q := sq.Select(
		"comments.id",
		"comments.body",
		"comments.post_id",
		"comments.user_id",
		"comments.updated_at",
		"comments.created_at",
	).
		From("public.comments").
		Limit(pageSize)
	// TODO: add filtering
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	if filter.PageSize != nil {
		q = q.Limit(*filter.PageSize)
	}
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &comments, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return comments, nil
}

func (r *CommentRepository) Update(
	ctx context.Context,
	comment *models.Comment,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.comments").
		Where(sq.Eq{"id": comment.ID}).
		Set("comments.body", comment.Body).
		Set("comments.post_id", comment.PostId).
		Set("comments.user_id", comment.UserId).
		Set("updated_at", comment.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("comment_id", fmt.Sprint(comment.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("comment_id", fmt.Sprint(comment.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("comment_id", fmt.Sprint(comment.ID))
		return e
	}
	return nil
}

func (r *CommentRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.comments").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("comment_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("comment_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("comment_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *CommentRepository) Count(
	ctx context.Context,
	filter *models.CommentFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.comments")
	// TODO: add filtering
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
