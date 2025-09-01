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

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewPostRepository(t *testing.T) {
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
		want  *PostRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				writeDB: mockDB,
				readDB:  mockDB,
			},
			want: &PostRepository{
				writeDB: mockDB,
				readDB:  mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewPostRepository(tt.args.readDB, tt.args.writeDB, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostRepository_Create(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.posts (id,created_at,updated_at,body) VALUES ($1,$2,$3,$4)"
	post := entities.NewMockPost(t)
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx  context.Context
		post entities.Post
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
						post.ID,
						post.UpdatedAt,
						post.CreatedAt,
						post.Body,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.ID,
						post.UpdatedAt,
						post.CreatedAt,
						post.Body,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.post)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostRepository_Get(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT posts.id, posts.created_at, posts.updated_at, posts.body FROM public.posts WHERE id = $1 LIMIT 1"
	post := entities.NewMockPost(t)
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
		want    entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newPostRows(t, []entities.Post{post})
				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnRows(rows)
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			want: entities.Post{},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", post.ID.String()),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			want:    entities.Post{},
			wantErr: errs.NewEntityNotFoundError().WithParam("post_id", post.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
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

func TestPostRepository_List(t *testing.T) {
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
	var listPosts []entities.Post
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listPosts = append(listPosts, entities.NewMockPost(t))
	}
	filter := entities.PostFilter{
		PageSize:   pointer.Of(uint64(10)),
		PageNumber: pointer.Of(uint64(2)),
		Search:     nil,
		OrderBy:    []entities.PostOrdering{"id"},
	}
	query := "SELECT posts.id, posts.created_at, posts.updated_at, posts.body FROM public.posts ORDER BY posts.id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx    context.Context
		filter entities.PostFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newPostRows(t, listPosts))
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
			want:    listPosts,
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
			r := &PostRepository{
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

func TestPostRepository_Update(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	post := entities.NewMockPost(t)
	query := `UPDATE public.posts SET created_at = $1, updated_at = $2, body = $3 WHERE id = $4`
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx  context.Context
		post entities.Post
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
						post.CreatedAt,
						post.UpdatedAt,
						post.Body,
						post.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Body,
						post.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("post_id", post.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Body,
						post.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", post.ID.String()),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Body,
						post.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", post.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Body,
						post.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				post: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", post.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.post)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostRepository_Delete(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	post := entities.NewMockPost(t)
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
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: nil,
		},
		{
			name: "post not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("post_id", post.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", post.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", post.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostRepository_Count(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	query := "SELECT count(id) FROM public.posts"
	ctx := context.Background()
	filter := entities.PostFilter{}
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx    context.Context
		filter entities.PostFilter
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
			r := &PostRepository{
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

func newPostRows(t *testing.T, listPosts []entities.Post) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"body",
		"updated_at",
		"created_at",
	})
	for _, post := range listPosts {
		rows.AddRow(
			post.ID,
			post.Body,
			post.UpdatedAt,
			post.CreatedAt,
		)
	}
	return rows
}
