package postgres

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	mock_models "github.com/mikalai-mitsin/example/internal/app/arch/models/mock"
	mock_postgres "github.com/mikalai-mitsin/example/internal/app/arch/repositories/postgres/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"go.uber.org/mock/gomock"

	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/arch/models"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewArchRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		database *sqlx.DB
		logger   Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *ArchRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &ArchRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewArchRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewArchRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_postgres.NewMockLogger(ctrl)
	query := "INSERT INTO public.arches (created_at,updated_at,name,title,subtitle,tags,versions,old_versions,release,tested,mark,submarine,numb) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id"
	arch := mock_models.NewArch(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Arch
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						arch.UpdatedAt,
						arch.CreatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(arch.ID, arch.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						arch.UpdatedAt,
						arch.CreatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArchRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArchRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_postgres.NewMockLogger(ctrl)
	query := "SELECT arches.id, arches.created_at, arches.updated_at, arches.name, arches.title, arches.subtitle, arches.tags, arches.versions, arches.old_versions, arches.release, arches.tested, arches.mark, arches.submarine, arches.numb FROM public.arches WHERE id = $1 LIMIT 1"
	arch := mock_models.NewArch(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   Logger
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newArchRows(t, []*models.Arch{arch})
				mock.ExpectQuery(query).WithArgs(arch.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			want:    arch,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(arch.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  arch.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("arch_id", string(arch.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(arch.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  arch.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError().WithParam("arch_id", string(arch.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArchRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchRepository_List(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_postgres.NewMockLogger(ctrl)
	ctx := context.Background()
	var listArches []*models.Arch
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listArches = append(listArches, mock_models.NewArch(t))
	}
	filter := &models.ArchFilter{
		PageSize:   pointer.Pointer(uint64(10)),
		PageNumber: pointer.Pointer(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
		IDs:        nil,
	}
	query := "SELECT arches.id, arches.created_at, arches.updated_at, arches.name, arches.title, arches.subtitle, arches.tags, arches.versions, arches.old_versions, arches.release, arches.tested, arches.mark, arches.submarine, arches.numb FROM public.arches ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database *sqlx.DB
		logger   Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ArchFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newArchRows(t, listArches))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listArches,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "error", Value: "test error"}},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArchRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_postgres.NewMockLogger(ctrl)
	arch := mock_models.NewArch(t)
	query := `UPDATE public.arches SET arches.created_at = $1, arches.updated_at = $2, arches.name = $3, arches.title = $4, arches.subtitle = $5, arches.tags = $6, arches.versions = $7, arches.old_versions = $8, arches.release = $9, arches.tested = $10, arches.mark = $11, arches.submarine = $12, arches.numb = $13 WHERE id = $14`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Arch
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						arch.CreatedAt,
						arch.UpdatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
						arch.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						arch.CreatedAt,
						arch.UpdatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
						arch.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("arch_id", string(arch.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						arch.CreatedAt,
						arch.UpdatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
						arch.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("arch_id", string(arch.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						arch.CreatedAt,
						arch.UpdatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
						arch.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("arch_id", string(arch.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						arch.CreatedAt,
						arch.UpdatedAt,
						arch.Name,
						arch.Title,
						arch.Subtitle,
						pq.Array(arch.Tags),
						pq.Array(arch.Versions),
						pq.Array(arch.OldVersions),
						arch.Release,
						arch.Tested,
						arch.Mark,
						arch.Submarine,
						arch.Numb,
						arch.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: arch,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("arch_id", string(arch.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArchRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArchRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_postgres.NewMockLogger(ctrl)
	arch := mock_models.NewArch(t)
	type fields struct {
		database *sqlx.DB
		logger   Logger
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				database: db,
				logger:   logger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.arches WHERE id = $1").
					WithArgs(arch.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  arch.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.arches WHERE id = $1").
					WithArgs(arch.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  arch.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("arch_id", string(arch.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.arches WHERE id = $1").
					WithArgs(arch.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  arch.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("arch_id", string(arch.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.arches WHERE id = $1").
					WithArgs(arch.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  arch.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("arch_id", string(arch.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArchRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArchRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT count(id) FROM public.arches"
	ctx := context.Background()
	filter := &models.ArchFilter{}
	type fields struct {
		database *sqlx.DB
		logger   Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ArchFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(1))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "bad return type",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow("one"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want: 0,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: errs.Params{
					{
						Key:   "error",
						Value: "sql: Scan error on column index 0, name \"count\": converting driver.Value type string (\"one\") to a uint64: invalid syntax",
					},
				},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    0,
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArchRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Count() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newArchRows(t *testing.T, listArches []*models.Arch) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"title",
		"subtitle",
		"tags",
		"versions",
		"old_versions",
		"release",
		"tested",
		"mark",
		"submarine",
		"numb",
		"updated_at",
		"created_at",
	})
	for _, arch := range listArches {
		rows.AddRow(
			arch.ID,
			arch.Name,
			arch.Title,
			arch.Subtitle,
			pq.Array(arch.Tags),
			pq.Array(arch.Versions),
			pq.Array(arch.OldVersions),
			arch.Release,
			arch.Tested,
			arch.Mark,
			arch.Submarine,
			arch.Numb,
			arch.UpdatedAt,
			arch.CreatedAt,
		)
	}
	return rows
}
