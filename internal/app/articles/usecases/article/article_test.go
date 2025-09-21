package usecases

import (
	"context"
	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewArticleUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := NewMockarticleService(ctrl)
	mockArticleEventService := NewMockarticleEventService(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		articleService      articleService
		articleEventService articleEventService
		dtxManager          dtxManager
		logger              logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *ArticleUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			want: &ArticleUseCase{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewArticleUseCase(
				tt.args.articleService,
				tt.args.articleEventService,
				tt.args.dtxManager,
				tt.args.logger,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := NewMockarticleService(ctrl)
	mockArticleEventService := NewMockarticleEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	type fields struct {
		articleService      articleService
		articleEventService articleEventService
		dtxManager          dtxManager
		logger              logger
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
				mockArticleService.EXPECT().
					Get(ctx, article.ID).
					Return(article, nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(article.ID),
			},
			want:    article,
			wantErr: nil,
		},
		{
			name: "Article not found",
			setup: func() {
				mockArticleService.EXPECT().
					Get(ctx, article.ID).
					Return(entities.Article{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
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
			i := &ArticleUseCase{
				articleService:      tt.fields.articleService,
				articleEventService: tt.fields.articleEventService,
				dtxManager:          tt.fields.dtxManager,
				logger:              tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := NewMockarticleService(ctrl)
	mockArticleEventService := NewMockarticleEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	create := entities.NewMockArticleCreate(t)
	type fields struct {
		articleService      articleService
		articleEventService articleEventService
		dtxManager          dtxManager
		logger              logger
	}
	type args struct {
		ctx    context.Context
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockArticleService.EXPECT().Create(ctx, mockTx, create).Return(article, nil)
				mockArticleEventService.EXPECT().Send(ctx, mockTx, article).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    article,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockArticleService.EXPECT().
					Create(ctx, mockTx, create).
					Return(entities.Article{}, errs.NewUnexpectedBehaviorError("c u"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Article{},
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ArticleUseCase{
				articleService:      tt.fields.articleService,
				articleEventService: tt.fields.articleEventService,
				dtxManager:          tt.fields.dtxManager,
				logger:              tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := NewMockarticleService(ctrl)
	mockArticleEventService := NewMockarticleEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	update := entities.NewMockArticleUpdate(t)
	type fields struct {
		articleService      articleService
		articleEventService articleEventService
		dtxManager          dtxManager
		logger              logger
	}
	type args struct {
		ctx    context.Context
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockArticleService.EXPECT().Update(ctx, mockTx, update).Return(article, nil)
				mockArticleEventService.EXPECT().Send(ctx, mockTx, article).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    article,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockArticleService.EXPECT().
					Update(ctx, mockTx, update).
					Return(entities.Article{}, errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Article{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ArticleUseCase{
				articleService:      tt.fields.articleService,
				articleEventService: tt.fields.articleEventService,
				dtxManager:          tt.fields.dtxManager,
				logger:              tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := NewMockarticleService(ctrl)
	mockArticleEventService := NewMockarticleEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	del := entities.NewMockArticleDelete(t)
	del.ID = article.ID
	type fields struct {
		articleService      articleService
		articleEventService articleEventService
		dtxManager          dtxManager
		logger              logger
	}
	type args struct {
		ctx context.Context
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockArticleService.EXPECT().
					Delete(ctx, mockTx, del).
					Return(article, nil)
				mockArticleEventService.EXPECT().Send(ctx, mockTx, article).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx: ctx,
				del: del,
			},
			want:    article,
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockArticleService.EXPECT().
					Delete(ctx, mockTx, del).
					Return(entities.Article{}, errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx: ctx,
				del: del,
			},
			want:    entities.Article{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ArticleUseCase{
				articleService:      tt.fields.articleService,
				articleEventService: tt.fields.articleEventService,
				dtxManager:          tt.fields.dtxManager,
				logger:              tt.fields.logger,
			}
			got, err := i.Delete(tt.args.ctx, tt.args.del)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleService := NewMockarticleService(ctrl)
	mockArticleEventService := NewMockarticleEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	filter := entities.NewMockArticleFilter(t)
	count := faker.New().UInt64Between(2, 20)
	articles := make([]entities.Article, 0, count)
	for i := uint64(0); i < count; i++ {
		articles = append(articles, entities.NewMockArticle(t))
	}
	type fields struct {
		articleService      articleService
		articleEventService articleEventService
		dtxManager          dtxManager
		logger              logger
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
				mockArticleService.EXPECT().
					List(ctx, filter).
					Return(articles, count, nil)
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
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
				mockArticleService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				articleService:      mockArticleService,
				articleEventService: mockArticleEventService,
				dtxManager:          mockDtxManager,
				logger:              mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ArticleUseCase{
				articleService:      tt.fields.articleService,
				articleEventService: tt.fields.articleEventService,
				dtxManager:          tt.fields.dtxManager,
				logger:              tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
