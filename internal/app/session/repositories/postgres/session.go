package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/session/models"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type SessionRepository struct {
	database *sqlx.DB
	logger   Logger
}

func NewSessionRepository(database *sqlx.DB, logger Logger) *SessionRepository {
	return &SessionRepository{database: database, logger: logger}
}

type SessionDTO struct {
	ID          string    `db:"id,omitempty"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}
type SessionListDTO []*SessionDTO

func (list SessionListDTO) ToModels() []*models.Session {
	items := make([]*models.Session, len(list))
	for i := range list {
		items[i] = list[i].ToModel()
	}
	return items
}
func NewSessionDTOFromModel(model *models.Session) *SessionDTO {
	dto := &SessionDTO{
		ID:          string(model.ID),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Title:       model.Title,
		Description: model.Description,
	}
	return dto
}
func (dto *SessionDTO) ToModel() *models.Session {
	model := &models.Session{
		ID:          uuid.UUID(dto.ID),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		Title:       dto.Title,
		Description: dto.Description,
	}
	return model
}
func (r *SessionRepository) Create(ctx context.Context, model *models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewSessionDTOFromModel(model)
	q := sq.Insert("public.sessions").
		Columns("created_at", "updated_at", "title", "description").
		Values(dto.CreatedAt, dto.UpdatedAt, dto.Title, dto.Description).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	model.ID = uuid.UUID(dto.ID)
	return nil
}

func (r *SessionRepository) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto SessionListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("sessions.id", "sessions.created_at", "sessions.updated_at", "sessions.title", "sessions.description").
		From("public.sessions").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"description"},
			},
		)
	}
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
	return dto.ToModels(), nil
}

func (r *SessionRepository) Count(
	ctx context.Context,
	filter *models.SessionFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.sessions")
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"description"},
			},
		)
	}
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
func (r *SessionRepository) Get(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &SessionDTO{}
	q := sq.Select("sessions.id", "sessions.created_at", "sessions.updated_at", "sessions.title", "sessions.description").
		From("public.sessions").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}
func (r *SessionRepository) Update(ctx context.Context, model *models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewSessionDTOFromModel(model)
	q := sq.Update("public.sessions").Where(sq.Eq{"id": model.ID})
	{
		q = q.Set("sessions.created_at", dto.CreatedAt)
		q = q.Set("sessions.updated_at", dto.UpdatedAt)
		q = q.Set("sessions.title", dto.Title)
		q = q.Set("sessions.description", dto.Description)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(model.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(model.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("session_id", fmt.Sprint(model.ID))
		return e
	}
	return nil
}
func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.sessions").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("session_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("session_id", fmt.Sprint(id))
		return e
	}
	return nil
}
