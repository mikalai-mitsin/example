package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type UserRepository struct {
	database *sqlx.DB
	logger   logger
}

func NewUserRepository(database *sqlx.DB, logger logger) *UserRepository {
	return &UserRepository{database: database, logger: logger}
}

type UserDTO struct {
	ID        string    `db:"id,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	GroupID   string    `db:"group_id"`
}
type UserListDTO []UserDTO

func (list UserListDTO) ToEntities() []entities.User {
	items := make([]entities.User, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewUserDTOFromEntity(entity entities.User) UserDTO {
	dto := UserDTO{
		ID:        string(entity.ID),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Password:  entity.Password,
		Email:     entity.Email,
		GroupID:   string(entity.GroupID),
	}
	return dto
}
func (dto UserDTO) toEntity() entities.User {
	entity := entities.User{
		ID:        uuid.UUID(dto.ID),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  dto.Password,
		Email:     dto.Email,
		GroupID:   entities.GroupID(dto.GroupID),
	}
	return entity
}
func (r *UserRepository) Create(ctx context.Context, entity entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewUserDTOFromEntity(entity)
	q := sq.Insert("public.users").
		Columns("created_at", "updated_at", "first_name", "last_name", "password", "email", "group_id").
		Values(dto.CreatedAt, dto.UpdatedAt, dto.FirstName, dto.LastName, dto.Password, dto.Email, dto.GroupID)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &UserDTO{}
	q := sq.Select("users.id", "users.created_at", "users.updated_at", "users.first_name", "users.last_name", "users.password", "users.email", "users.group_id").
		From("public.users").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("user_id", string(id))
		return entities.User{}, e
	}
	return dto.toEntity(), nil
}

func (r *UserRepository) List(
	ctx context.Context,
	filter entities.UserFilter,
) ([]entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto UserListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("users.id", "users.created_at", "users.updated_at", "users.first_name", "users.last_name", "users.password", "users.email", "users.group_id").
		From("public.users").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"first_name", "last_name", "email"},
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
	return dto.ToEntities(), nil
}
func (r *UserRepository) Count(ctx context.Context, filter entities.UserFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.users")
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"first_name", "last_name", "email"},
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
func (r *UserRepository) Update(ctx context.Context, entity entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewUserDTOFromEntity(entity)
	q := sq.Update("public.users").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("first_name", dto.FirstName)
		q = q.Set("last_name", dto.LastName)
		q = q.Set("password", dto.Password)
		q = q.Set("email", dto.Email)
		q = q.Set("group_id", dto.GroupID)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("user_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("user_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("user_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.users").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("user_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("user_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("user_id", fmt.Sprint(id))
		return e
	}
	return nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &UserDTO{}
	q := sq.Select("users.id", "users.created_at", "users.updated_at", "users.first_name", "users.last_name", "users.password", "users.email", "users.group_id").
		From("public.users").
		Where(sq.Eq{"email": email}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("user_email", email)
		return entities.User{}, e
	}
	return dto.toEntity(), nil
}
