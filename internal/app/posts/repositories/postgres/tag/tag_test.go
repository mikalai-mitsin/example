package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewTagRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		writeDB database
		readDB  database
		logger  logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *TagRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				writeDB: mockDB,
				readDB:  mockDB,
			},
			want: &TagRepository{
				writeDB: mockDB,
				readDB:  mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewTagRepository(tt.args.readDB, tt.args.writeDB, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagRepository_Create(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.tags (id,created_at,updated_at,post_id,value) VALUES ($1,$2,$3,$4,$5)"
	tag := entities.NewMockTag(t)
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx context.Context
		tag entities.Tag
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
						tag.ID,
						tag.UpdatedAt,
						tag.CreatedAt,
						tag.PostId,
						tag.Value,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						tag.ID,
						tag.UpdatedAt,
						tag.CreatedAt,
						tag.PostId,
						tag.Value,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &TagRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.tag)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTagRepository_Get(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT tags.id, tags.created_at, tags.updated_at, tags.post_id, tags.value FROM public.tags WHERE id = $1 LIMIT 1"
	tag := entities.NewMockTag(t)
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
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
		want    entities.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newTagRows(t, []entities.Tag{tag})
				mock.ExpectQuery(query).WithArgs(tag.ID).WillReturnRows(rows)
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  tag.ID,
			},
			want:    tag,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(tag.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  tag.ID,
			},
			want: entities.Tag{},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("tag_id", tag.ID.String()),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(tag.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  tag.ID,
			},
			want:    entities.Tag{},
			wantErr: errs.NewEntityNotFoundError().WithParam("tag_id", tag.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &TagRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagRepository_List(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listTags []entities.Tag
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listTags = append(listTags, entities.NewMockTag(t))
	}
	filter := entities.TagFilter{
		PageSize:   pointer.Of(uint64(10)),
		PageNumber: pointer.Of(uint64(2)),
		Search:     nil,
		OrderBy:    []entities.TagOrdering{"id"},
	}
	query := "SELECT tags.id, tags.created_at, tags.updated_at, tags.post_id, tags.value FROM public.tags ORDER BY tags.id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx    context.Context
		filter entities.TagFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newTagRows(t, listTags))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listTags,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
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
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
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
			r := &TagRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagRepository_Update(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	tag := entities.NewMockTag(t)
	query := `UPDATE public.tags SET created_at = $1, updated_at = $2, post_id = $3, value = $4 WHERE id = $5`
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx context.Context
		tag entities.Tag
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
						tag.CreatedAt,
						tag.UpdatedAt,
						tag.PostId,
						tag.Value,
						tag.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						tag.CreatedAt,
						tag.UpdatedAt,
						tag.PostId,
						tag.Value,
						tag.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("tag_id", tag.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						tag.CreatedAt,
						tag.UpdatedAt,
						tag.PostId,
						tag.Value,
						tag.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("tag_id", tag.ID.String()),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						tag.CreatedAt,
						tag.UpdatedAt,
						tag.PostId,
						tag.Value,
						tag.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("tag_id", tag.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						tag.CreatedAt,
						tag.UpdatedAt,
						tag.PostId,
						tag.Value,
						tag.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				tag: tag,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("tag_id", tag.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &TagRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.tag)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTagRepository_Delete(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	tag := entities.NewMockTag(t)
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
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
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.tags WHERE id = $1").
					WithArgs(tag.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  tag.ID,
			},
			wantErr: nil,
		},
		{
			name: "tag not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.tags WHERE id = $1").
					WithArgs(tag.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  tag.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("tag_id", tag.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.tags WHERE id = $1").
					WithArgs(tag.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  tag.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("tag_id", tag.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.tags WHERE id = $1").
					WithArgs(tag.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  tag.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("tag_id", tag.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &TagRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTagRepository_Count(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	query := "SELECT count(id) FROM public.tags"
	ctx := context.Background()
	filter := entities.TagFilter{}
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx    context.Context
		filter entities.TagFilter
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
				writeDB: mockDB,
				readDB:  mockDB,
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
				writeDB: mockDB,
				readDB:  mockDB,
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
				writeDB: mockDB,
				readDB:  mockDB,
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
			r := &TagRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newTagRows(t *testing.T, listTags []entities.Tag) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"post_id",
		"value",
		"updated_at",
		"created_at",
	})
	for _, tag := range listTags {
		rows.AddRow(
			tag.ID,
			tag.PostId,
			tag.Value,
			tag.UpdatedAt,
			tag.CreatedAt,
		)
	}
	return rows
}
