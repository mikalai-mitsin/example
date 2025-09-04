package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewArticleServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleUseCase := NewMockarticleUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		articleUseCase articleUseCase
		logger         logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.ArticleServiceServer
	}{
		{
			name: "ok",
			args: args{
				articleUseCase: mockArticleUseCase,
				logger:         mockLogger,
			},
			want: &ArticleServiceServer{
				articleUseCase: mockArticleUseCase,
				logger:         mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArticleServiceServer(tt.args.articleUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleUseCase := NewMockarticleUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := entities.NewMockArticleCreate(t)
	article := entities.NewMockArticle(t)
	type fields struct {
		UnimplementedArticleServiceServer examplepb.UnimplementedArticleServiceServer
		articleUseCase                    articleUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArticleCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Article
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockArticleUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(article, nil)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.ArticleCreate{},
			},
			want:    decodeArticle(article),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockArticleUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(entities.Article{}, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.ArticleCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := ArticleServiceServer{
				UnimplementedArticleServiceServer: tt.fields.UnimplementedArticleServiceServer,
				articleUseCase:                    tt.fields.articleUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleUseCase := NewMockarticleUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedArticleServiceServer examplepb.UnimplementedArticleServiceServer
		articleUseCase                    articleUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArticleDelete
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockArticleUseCase.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArticleDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockArticleUseCase.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArticleDelete{
					Id: id.String(),
				},
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "details", Value: "i error"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := ArticleServiceServer{
				UnimplementedArticleServiceServer: tt.fields.UnimplementedArticleServiceServer,
				articleUseCase:                    tt.fields.articleUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleUseCase := NewMockarticleUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	type fields struct {
		UnimplementedArticleServiceServer examplepb.UnimplementedArticleServiceServer
		articleUseCase                    articleUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArticleGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Article
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockArticleUseCase.EXPECT().Get(ctx, article.ID).Return(article, nil).Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArticleGet{
					Id: article.ID.String(),
				},
			},
			want:    decodeArticle(article),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockArticleUseCase.EXPECT().Get(ctx, article.ID).
					Return(entities.Article{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArticleGet{
					Id: article.ID.String(),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := ArticleServiceServer{
				UnimplementedArticleServiceServer: tt.fields.UnimplementedArticleServiceServer,
				articleUseCase:                    tt.fields.articleUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleUseCase := NewMockarticleUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := entities.NewMockArticleFilter(t)
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListArticle{
		Items: make([]*examplepb.Article, 0, int(count)),
		Count: count,
	}
	listArticles := make([]entities.Article, 0, int(count))
	type fields struct {
		UnimplementedArticleServiceServer examplepb.UnimplementedArticleServiceServer
		articleUseCase                    articleUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArticleFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListArticle
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockArticleUseCase.EXPECT().
					List(ctx, gomock.Any()).
					Return(listArticles, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArticleFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    nil,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockArticleUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.ArticleFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    nil,
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := ArticleServiceServer{
				UnimplementedArticleServiceServer: tt.fields.UnimplementedArticleServiceServer,
				articleUseCase:                    tt.fields.articleUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArticleUseCase := NewMockarticleUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	update := entities.NewMockArticleUpdate(t)
	type fields struct {
		UnimplementedArticleServiceServer examplepb.UnimplementedArticleServiceServer
		articleUseCase                    articleUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.ArticleUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Article
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockArticleUseCase.EXPECT().Update(ctx, gomock.Any()).Return(article, nil).Times(1)
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeArticleUpdate(update),
			},
			want:    decodeArticle(article),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockArticleUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(entities.Article{}, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedArticleServiceServer: examplepb.UnimplementedArticleServiceServer{},
				articleUseCase:                    mockArticleUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeArticleUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := ArticleServiceServer{
				UnimplementedArticleServiceServer: tt.fields.UnimplementedArticleServiceServer,
				articleUseCase:                    tt.fields.articleUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeArticle(t *testing.T) {
	article := entities.NewMockArticle(t)
	result := &examplepb.Article{
		Id:          article.ID.String(),
		UpdatedAt:   timestamppb.New(article.UpdatedAt),
		CreatedAt:   timestamppb.New(article.CreatedAt),
		Title:       string(article.Title),
		Subtitle:    string(article.Subtitle),
		Body:        string(article.Body),
		IsPublished: bool(article.IsPublished),
	}
	type args struct {
		article entities.Article
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Article
	}{
		{
			name: "ok",
			args: args{
				article: article,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeArticle(tt.args.article)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeArticleFilter(t *testing.T) {
	type args struct {
		input *examplepb.ArticleFilter
	}
	tests := []struct {
		name string
		args args
		want entities.ArticleFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.ArticleFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
				},
			},
			want: entities.ArticleFilter{
				PageSize:   pointer.Of(uint64(5)),
				PageNumber: pointer.Of(uint64(2)),
				OrderBy:    []entities.ArticleOrdering{"created_at", "id"},
				Search:     pointer.Of("my name is"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeArticleFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
