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

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewArticleRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		database database
		logger   logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *ArticleRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &ArticleRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewArticleRepository(tt.args.database, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleRepository_Create(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.articles (id,created_at,updated_at,title,subtitle,body,is_published) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	article := entities.NewMockArticle(t)
	ctx := context.Background()
	type fields struct {
		database database
		logger   logger
	}
	type args struct {
		ctx     context.Context
		article entities.Article
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
						article.ID,
						article.UpdatedAt,
						article.CreatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						article.ID,
						article.UpdatedAt,
						article.CreatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArticleRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.article)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestArticleRepository_Get(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT articles.id, articles.created_at, articles.updated_at, articles.title, articles.subtitle, articles.body, articles.is_published FROM public.articles WHERE id = $1 LIMIT 1"
	article := entities.NewMockArticle(t)
	ctx := context.Background()
	type fields struct {
		database database
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
		want    entities.Article
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newArticleRows(t, []entities.Article{article})
				mock.ExpectQuery(query).WithArgs(article.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  article.ID,
			},
			want:    article,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(article.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  article.ID,
			},
			want: entities.Article{},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("article_id", article.ID.String()),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(article.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  article.ID,
			},
			want:    entities.Article{},
			wantErr: errs.NewEntityNotFoundError().WithParam("article_id", article.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArticleRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleRepository_List(t *testing.T) {
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
	var listArticles []entities.Article
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listArticles = append(listArticles, entities.NewMockArticle(t))
	}
	filter := entities.ArticleFilter{
		PageSize:   pointer.Of(uint64(10)),
		PageNumber: pointer.Of(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
	}
	query := "SELECT articles.id, articles.created_at, articles.updated_at, articles.title, articles.subtitle, articles.body, articles.is_published FROM public.articles ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database database
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter entities.ArticleFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Article
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newArticleRows(t, listArticles))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listArticles,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: mockDB,
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
				database: mockDB,
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
			r := &ArticleRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleRepository_Update(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	article := entities.NewMockArticle(t)
	query := `UPDATE public.articles SET created_at = $1, updated_at = $2, title = $3, subtitle = $4, body = $5, is_published = $6 WHERE id = $7`
	ctx := context.Background()
	type fields struct {
		database database
		logger   logger
	}
	type args struct {
		ctx     context.Context
		article entities.Article
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
						article.CreatedAt,
						article.UpdatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
						article.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						article.CreatedAt,
						article.UpdatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
						article.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("article_id", article.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						article.CreatedAt,
						article.UpdatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
						article.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("article_id", article.ID.String()),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						article.CreatedAt,
						article.UpdatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
						article.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("article_id", article.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						article.CreatedAt,
						article.UpdatedAt,
						article.Title,
						article.Subtitle,
						article.Body,
						article.IsPublished,
						article.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				article: article,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("article_id", article.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArticleRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.article)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestArticleRepository_Delete(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	article := entities.NewMockArticle(t)
	type fields struct {
		database database
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
				database: mockDB,
				logger:   mockLogger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.articles WHERE id = $1").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  article.ID,
			},
			wantErr: nil,
		},
		{
			name: "article not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.articles WHERE id = $1").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  article.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("article_id", article.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.articles WHERE id = $1").
					WithArgs(article.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  article.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("article_id", article.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.articles WHERE id = $1").
					WithArgs(article.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: mockDB,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  article.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("article_id", article.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &ArticleRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestArticleRepository_Count(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	query := "SELECT count(id) FROM public.articles"
	ctx := context.Background()
	filter := entities.ArticleFilter{}
	type fields struct {
		database database
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter entities.ArticleFilter
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
				database: mockDB,
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
				database: mockDB,
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
				database: mockDB,
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
			r := &ArticleRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newArticleRows(t *testing.T, listArticles []entities.Article) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"title",
		"subtitle",
		"body",
		"is_published",
		"updated_at",
		"created_at",
	})
	for _, article := range listArticles {
		rows.AddRow(
			article.ID,
			article.Title,
			article.Subtitle,
			article.Body,
			article.IsPublished,
			article.UpdatedAt,
			article.CreatedAt,
		)
	}
	return rows
}
