package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewArticleService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleRepository := NewMockarticleRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	type args struct {
		articleRepository articleRepository
		clock             clock
		logger            logger
		uuid              uuidGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *ArticleService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
			want: &ArticleService{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewArticleService(
				tt.args.articleRepository,
				tt.args.clock,
				tt.args.logger,
				tt.args.uuid,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleRepository := NewMockarticleRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	type fields struct {
		articleRepository articleRepository
		logger            logger
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
				mockArticleRepository.EXPECT().Get(ctx, article.ID).Return(article, nil)
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  article.ID,
			},
			want:    article,
			wantErr: nil,
		},
		{
			name: "Article not found",
			setup: func() {
				mockArticleRepository.EXPECT().
					Get(ctx, article.ID).
					Return(entities.Article{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  article.ID,
			},
			want:    entities.Article{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArticleService{
				articleRepository: tt.fields.articleRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleRepository := NewMockarticleRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var articles []entities.Article
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		articles = append(articles, entities.NewMockArticle(t))
	}
	filter := entities.NewMockArticleFilter(t)
	type fields struct {
		articleRepository articleRepository
		logger            logger
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
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockArticleRepository.EXPECT().List(ctx, filter).Return(articles, nil)
				mockArticleRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    articles,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockArticleRepository.EXPECT().
					List(ctx, filter).
					Return([]entities.Article{}, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "count error",
			setup: func() {
				mockArticleRepository.EXPECT().List(ctx, filter).Return(articles, nil)
				mockArticleRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArticleService{
				articleRepository: tt.fields.articleRepository,
				logger:            tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestArticleService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleRepository := NewMockarticleRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	create := entities.NewMockArticleCreate(t)
	now := time.Now().UTC()
	type fields struct {
		articleRepository articleRepository
		clock             clock
		logger            logger
		uuid              uuidGenerator
	}
	type args struct {
		ctx    context.Context
		tx     dtx.TX
		create entities.ArticleCreate
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
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().
					NewUUID().
					Return(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
				mockArticleRepository.EXPECT().
					Create(
						ctx,
						mockTx,
						entities.Article{
							ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							Title:       create.Title,
							Subtitle:    create.Subtitle,
							Body:        create.Body,
							IsPublished: create.IsPublished,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(nil)
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: create,
			},
			want: entities.Article{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:       create.Title,
				Subtitle:    create.Subtitle,
				Body:        create.Body,
				IsPublished: create.IsPublished,
				UpdatedAt:   now,
				CreatedAt:   now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().
					NewUUID().
					Return(uuid.MustParse("00000000-0000-0000-0000-000000000002"))
				mockArticleRepository.EXPECT().
					Create(
						ctx,
						mockTx,
						entities.Article{
							ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							Title:       create.Title,
							Subtitle:    create.Subtitle,
							Body:        create.Body,
							IsPublished: create.IsPublished,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: create,
			},
			want:    entities.Article{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
				clock:             mockClock,
				uuid:              mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: entities.ArticleCreate{},
			},
			want: entities.Article{},
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "title", Value: "cannot be blank"},
				errs.Param{Key: "subtitle", Value: "cannot be blank"},
				errs.Param{Key: "body", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArticleService{
				articleRepository: tt.fields.articleRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
				uuid:              tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.tx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleRepository := NewMockarticleRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	mockClock := NewMockclock(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	update := entities.NewMockArticleUpdate(t)
	now := time.Now().UTC()
	updatedArticle := entities.Article{
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
		DeletedAt: article.DeletedAt,
		UpdatedAt: now,

		Title:       *update.Title,
		Subtitle:    *update.Subtitle,
		Body:        *update.Body,
		IsPublished: *update.IsPublished,
	}
	type fields struct {
		articleRepository articleRepository
		clock             clock
		logger            logger
	}
	type args struct {
		ctx    context.Context
		tx     dtx.TX
		update entities.ArticleUpdate
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
				mockClock.EXPECT().Now().Return(now)
				mockArticleRepository.EXPECT().
					Get(ctx, update.ID).Return(article, nil)
				mockArticleRepository.EXPECT().
					Update(ctx, mockTx, updatedArticle).Return(nil)
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    updatedArticle,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockArticleRepository.EXPECT().
					Get(ctx, update.ID).
					Return(article, nil)
				mockArticleRepository.EXPECT().
					Update(ctx, mockTx, updatedArticle).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    entities.Article{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Article not found",
			setup: func() {
				mockArticleRepository.EXPECT().
					Get(ctx, update.ID).
					Return(entities.Article{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    entities.Article{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArticleService{
				articleRepository: tt.fields.articleRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.tx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleRepository := NewMockarticleRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockClock := NewMockclock(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	now := time.Now().UTC()
	article := entities.NewMockArticle(t)
	deletedArticle := entities.Article{
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.CreatedAt,
		DeletedAt: pointer.Of(now),

		Title:       article.Title,
		Subtitle:    article.Subtitle,
		Body:        article.Body,
		IsPublished: article.IsPublished,
	}
	del := entities.NewMockArticleDelete(t)
	del.ID = article.ID
	type fields struct {
		articleRepository articleRepository
		clock             clock
		logger            logger
	}
	type args struct {
		ctx context.Context
		tx  dtx.TX
		del entities.ArticleDelete
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
				mockClock.EXPECT().Now().Return(now)
				mockArticleRepository.EXPECT().
					Get(ctx, del.ID).
					Return(article, nil)
				mockArticleRepository.EXPECT().
					Update(ctx, mockTx, deletedArticle).
					Return(nil)
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
				clock:             mockClock,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				del: del,
			},
			want:    deletedArticle,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockArticleRepository.EXPECT().
					Get(ctx, del.ID).
					Return(article, nil)
				mockArticleRepository.EXPECT().
					Update(ctx, mockTx, deletedArticle).
					Return(errs.NewUnexpectedBehaviorError("test error 12"))
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
				clock:             mockClock,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				del: del,
			},
			want:    entities.Article{},
			wantErr: errs.NewUnexpectedBehaviorError("test error 12"),
		},
		{
			name: "Article not found",
			setup: func() {
				mockArticleRepository.EXPECT().
					Get(ctx, del.ID).
					Return(entities.Article{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				articleRepository: mockArticleRepository,
				logger:            mockLogger,
				clock:             mockClock,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				del: del,
			},
			want:    entities.Article{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArticleService{
				articleRepository: tt.fields.articleRepository,
				logger:            tt.fields.logger,
				clock:             tt.fields.clock,
			}
			got, err := u.Delete(tt.args.ctx, tt.args.tx, tt.args.del)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
