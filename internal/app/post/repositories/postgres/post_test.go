package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/post/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/post/entities"
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
		database *sqlx.DB
		logger   logger
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
				database: mockDB,
			},
			want: &PostRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewPostRepository(tt.args.database, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.posts (created_at,updated_at,title,order,is_optional) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	post := mock_entities.NewPost(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx  context.Context
		card *entities.Post
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
						post.UpdatedAt,
						post.CreatedAt,
						post.Title,
						post.Order,
						post.IsOptional,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(post.ID, post.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						post.UpdatedAt,
						post.CreatedAt,
						post.Title,
						post.Order,
						post.IsOptional,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.card)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.order, posts.is_optional FROM public.posts WHERE id = $1 LIMIT 1"
	post := mock_entities.NewPost(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
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
		want    *entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newPostRows(t, []*entities.Post{post})
				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
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
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", string(post.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError().WithParam("post_id", string(post.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostRepository_List(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listPosts []*entities.Post
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listPosts = append(listPosts, mock_entities.NewPost(t))
	}
	filter := &entities.PostFilter{
		PageSize:   pointer.Pointer(uint64(10)),
		PageNumber: pointer.Pointer(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
		IDs:        nil,
	}
	query := "SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.order, posts.is_optional FROM public.posts ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.PostFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newPostRows(t, listPosts))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
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
				database: db,
				logger:   mockLogger,
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
				logger:   mockLogger,
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
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	post := mock_entities.NewPost(t)
	query := `UPDATE public.posts SET posts.created_at = $1, posts.updated_at = $2, posts.title = $3, posts.order = $4, posts.is_optional = $5 WHERE id = $6`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx  context.Context
		card *entities.Post
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
						post.Title,
						post.Order,
						post.IsOptional,
						post.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
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
						post.Title,
						post.Order,
						post.IsOptional,
						post.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("post_id", string(post.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Title,
						post.Order,
						post.IsOptional,
						post.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", string(post.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Title,
						post.Order,
						post.IsOptional,
						post.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", string(post.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						post.CreatedAt,
						post.UpdatedAt,
						post.Title,
						post.Order,
						post.IsOptional,
						post.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: post,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", string(post.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.card)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	post := mock_entities.NewPost(t)
	type fields struct {
		database *sqlx.DB
		logger   logger
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
				logger:   mockLogger,
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
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("post_id", string(post.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", string(post.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.posts WHERE id = $1").
					WithArgs(post.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  post.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("post_id", string(post.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &PostRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT count(id) FROM public.posts"
	ctx := context.Background()
	filter := &entities.PostFilter{}
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.PostFilter
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
			r := &PostRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newPostRows(t *testing.T, listPosts []*entities.Post) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"title",
		"order",
		"is_optional",
		"updated_at",
		"created_at",
	})
	for _, post := range listPosts {
		rows.AddRow(
			post.ID,
			post.Title,
			post.Order,
			post.IsOptional,
			post.UpdatedAt,
			post.CreatedAt,
		)
	}
	return rows
}
