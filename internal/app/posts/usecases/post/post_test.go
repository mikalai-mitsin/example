package usecases

import (
	"context"
	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewPostUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockPostEventService := NewMockpostEventService(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		postService      postService
		postEventService postEventService
		dtxManager       dtxManager
		logger           logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *PostUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			want: &PostUseCase{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewPostUseCase(
				tt.args.postService,
				tt.args.postEventService,
				tt.args.dtxManager,
				tt.args.logger,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockPostEventService := NewMockpostEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		postService      postService
		postEventService postEventService
		dtxManager       dtxManager
		logger           logger
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
				mockPostService.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(post.ID),
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "Post not found",
			setup: func() {
				mockPostService.EXPECT().
					Get(ctx, post.ID).
					Return(entities.Post{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			want:    entities.Post{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:      tt.fields.postService,
				postEventService: tt.fields.postEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockPostEventService := NewMockpostEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	create := entities.NewMockPostCreate(t)
	type fields struct {
		postService      postService
		postEventService postEventService
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
		create entities.PostCreate
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockPostService.EXPECT().Create(ctx, mockTx, create).Return(post, nil)
				mockPostEventService.EXPECT().Created(ctx, mockTx, post).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockPostService.EXPECT().
					Create(ctx, mockTx, create).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("c u"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Post{},
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:      tt.fields.postService,
				postEventService: tt.fields.postEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockPostEventService := NewMockpostEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	update := entities.NewMockPostUpdate(t)
	type fields struct {
		postService      postService
		postEventService postEventService
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
		update entities.PostUpdate
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockPostService.EXPECT().Update(ctx, mockTx, update).Return(post, nil)
				mockPostEventService.EXPECT().Updated(ctx, mockTx, post).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockPostService.EXPECT().
					Update(ctx, mockTx, update).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Post{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:      tt.fields.postService,
				postEventService: tt.fields.postEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockPostEventService := NewMockpostEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		postService      postService
		postEventService postEventService
		dtxManager       dtxManager
		logger           logger
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
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockPostService.EXPECT().
					Delete(ctx, mockTx, post.ID).
					Return(nil)
				mockPostEventService.EXPECT().Deleted(ctx, mockTx, post.ID).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockPostService.EXPECT().
					Delete(ctx, mockTx, post.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:      tt.fields.postService,
				postEventService: tt.fields.postEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockPostEventService := NewMockpostEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	filter := entities.NewMockPostFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listPosts := make([]entities.Post, 0, count)
	for i := uint64(0); i < count; i++ {
		listPosts = append(listPosts, entities.NewMockPost(t))
	}
	type fields struct {
		postService      postService
		postEventService postEventService
		dtxManager       dtxManager
		logger           logger
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
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostService.EXPECT().
					List(ctx, filter).
					Return(listPosts, count, nil)
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listPosts,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockPostService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				postService:      mockPostService,
				postEventService: mockPostEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
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
			i := &PostUseCase{
				postService:      tt.fields.postService,
				postEventService: tt.fields.postEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
