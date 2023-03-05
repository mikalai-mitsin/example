package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/postgresql"
	"github.com/018bf/example/pkg/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ArchRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewArchRepository(database *sqlx.DB, logger log.Logger) repositories.ArchRepository {
	return &ArchRepository{database: database, logger: logger}
}
func (r *ArchRepository) Create(ctx context.Context, arch *models.Arch) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewArchDTOFromModel(arch)
	q := sq.Insert("public.arches").
		Columns("updated_at", "created_at", "name", "title", "subtitle", "tags", "versions", "old_versions", "release", "tested", "mark", "submarine", "numb").
		Values(dto.UpdatedAt, dto.CreatedAt, dto.Name, dto.Title, dto.Subtitle, dto.Tags, dto.Versions, dto.OldVersions, dto.Release, dto.Tested, dto.Mark, dto.Submarine, dto.Numb).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	arch.ID = models.UUID(dto.ID)
	return nil
}
func (r *ArchRepository) Get(ctx context.Context, id models.UUID) (*models.Arch, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &ArchDTO{}
	q := sq.Select("arches.id", "arches.updated_at", "arches.created_at", "arches.name", "arches.title", "arches.subtitle", "arches.tags", "arches.versions", "arches.old_versions", "arches.release", "arches.tested", "arches.mark", "arches.submarine", "arches.numb").
		From("public.arches").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("arch_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *ArchRepository) List(
	ctx context.Context,
	filter *models.ArchFilter,
) ([]*models.Arch, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto ArchListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select("arches.id", "arches.updated_at", "arches.created_at", "arches.name", "arches.title", "arches.subtitle", "arches.tags", "arches.versions", "arches.old_versions", "arches.release", "arches.tested", "arches.mark", "arches.submarine", "arches.numb").
		From("public.arches").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"name", "subtitle"},
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
func (r *ArchRepository) Count(ctx context.Context, filter *models.ArchFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.arches")
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"name", "subtitle"},
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
func (r *ArchRepository) Update(ctx context.Context, arch *models.Arch) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewArchDTOFromModel(arch)
	q := sq.Update("public.arches").Where(sq.Eq{"id": arch.ID})
	{
		q = q.Set("arches.updated_at", dto.UpdatedAt)
		q = q.Set("arches.name", dto.Name)
		q = q.Set("arches.title", dto.Title)
		q = q.Set("arches.subtitle", dto.Subtitle)
		q = q.Set("arches.tags", dto.Tags)
		q = q.Set("arches.versions", dto.Versions)
		q = q.Set("arches.old_versions", dto.OldVersions)
		q = q.Set("arches.release", dto.Release)
		q = q.Set("arches.tested", dto.Tested)
		q = q.Set("arches.mark", dto.Mark)
		q = q.Set("arches.submarine", dto.Submarine)
		q = q.Set("arches.numb", dto.Numb)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("arch_id", fmt.Sprint(arch.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("arch_id", fmt.Sprint(arch.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("arch_id", fmt.Sprint(arch.ID))
		return e
	}
	return nil
}
func (r *ArchRepository) Delete(ctx context.Context, id models.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.arches").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("arch_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("arch_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("arch_id", fmt.Sprint(id))
		return e
	}
	return nil
}

type ArchDTO struct {
	ID          string         `db:"id,omitempty"`
	UpdatedAt   time.Time      `db:"updated_at,omitempty"`
	CreatedAt   time.Time      `db:"created_at,omitempty"`
	Name        string         `db:"name"`
	Title       string         `db:"title"`
	Subtitle    string         `db:"subtitle"`
	Tags        pq.StringArray `db:"tags"`
	Versions    pq.Int64Array  `db:"versions"`
	OldVersions pq.Int64Array  `db:"old_versions"`
	Release     time.Time      `db:"release"`
	Tested      time.Time      `db:"tested"`
	Mark        string         `db:"mark"`
	Submarine   string         `db:"submarine"`
	Numb        int64          `db:"numb"`
}
type ArchListDTO []*ArchDTO

func (list ArchListDTO) ToModels() []*models.Arch {
	listArches := make([]*models.Arch, len(list))
	for i := range list {
		listArches[i] = list[i].ToModel()
	}
	return listArches
}
func NewArchDTOFromModel(arch *models.Arch) *ArchDTO {
	dto := &ArchDTO{
		ID:          string(arch.ID),
		UpdatedAt:   arch.UpdatedAt,
		CreatedAt:   arch.CreatedAt,
		Name:        arch.Name,
		Title:       arch.Title,
		Subtitle:    arch.Subtitle,
		Tags:        pq.StringArray{},
		Versions:    pq.Int64Array{},
		OldVersions: pq.Int64Array{},
		Release:     arch.Release,
		Tested:      arch.Tested,
		Mark:        arch.Mark,
		Submarine:   arch.Submarine,
		Numb:        int64(arch.Numb),
	}
	for _, param := range arch.Tags {
		dto.Tags = append(dto.Tags, param)
	}
	for _, param := range arch.Versions {
		dto.Versions = append(dto.Versions, int64(param))
	}
	for _, param := range arch.OldVersions {
		dto.OldVersions = append(dto.OldVersions, int64(param))
	}
	return dto
}
func (dto *ArchDTO) ToModel() *models.Arch {
	model := &models.Arch{
		ID:          models.UUID(dto.ID),
		UpdatedAt:   dto.UpdatedAt,
		CreatedAt:   dto.CreatedAt,
		Name:        dto.Name,
		Title:       dto.Title,
		Subtitle:    dto.Subtitle,
		Tags:        []string{},
		Versions:    []uint{},
		OldVersions: []uint64{},
		Release:     dto.Release,
		Tested:      dto.Tested,
		Mark:        dto.Mark,
		Submarine:   dto.Submarine,
		Numb:        uint64(dto.Numb),
	}
	for _, param := range dto.Tags {
		model.Tags = append(model.Tags, param)
	}
	for _, param := range dto.Versions {
		model.Versions = append(model.Versions, uint(param))
	}
	for _, param := range dto.OldVersions {
		model.OldVersions = append(model.OldVersions, uint64(param))
	}
	return model
}
