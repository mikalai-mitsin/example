package usecases

import (
	"context"
	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewLikeUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeService := NewMocklikeService(ctrl)
	mockLikeEventService := NewMocklikeEventService(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		likeService      likeService
		likeEventService likeEventService
		dtxManager       dtxManager
		logger           logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *LikeUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			want: &LikeUseCase{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewLikeUseCase(
				tt.args.likeService,
				tt.args.likeEventService,
				tt.args.dtxManager,
				tt.args.logger,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeService := NewMocklikeService(ctrl)
	mockLikeEventService := NewMocklikeEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		likeService      likeService
		likeEventService likeEventService
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
		want    entities.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeService.EXPECT().
					Get(ctx, like.ID).
					Return(like, nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(like.ID),
			},
			want:    like,
			wantErr: nil,
		},
		{
			name: "Like not found",
			setup: func() {
				mockLikeService.EXPECT().
					Get(ctx, like.ID).
					Return(entities.Like{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			want:    entities.Like{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &LikeUseCase{
				likeService:      tt.fields.likeService,
				likeEventService: tt.fields.likeEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeService := NewMocklikeService(ctrl)
	mockLikeEventService := NewMocklikeEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	create := entities.NewMockLikeCreate(t)
	type fields struct {
		likeService      likeService
		likeEventService likeEventService
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
		create entities.LikeCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    entities.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockLikeService.EXPECT().Create(ctx, mockTx, create).Return(like, nil)
				mockLikeEventService.EXPECT().Created(ctx, mockTx, like).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    like,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockLikeService.EXPECT().
					Create(ctx, mockTx, create).
					Return(entities.Like{}, errs.NewUnexpectedBehaviorError("c u"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Like{},
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &LikeUseCase{
				likeService:      tt.fields.likeService,
				likeEventService: tt.fields.likeEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeService := NewMocklikeService(ctrl)
	mockLikeEventService := NewMocklikeEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	update := entities.NewMockLikeUpdate(t)
	type fields struct {
		likeService      likeService
		likeEventService likeEventService
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
		update entities.LikeUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    entities.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockLikeService.EXPECT().Update(ctx, mockTx, update).Return(like, nil)
				mockLikeEventService.EXPECT().Updated(ctx, mockTx, like).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    like,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockLikeService.EXPECT().
					Update(ctx, mockTx, update).
					Return(entities.Like{}, errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Like{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &LikeUseCase{
				likeService:      tt.fields.likeService,
				likeEventService: tt.fields.likeEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeService := NewMocklikeService(ctrl)
	mockLikeEventService := NewMocklikeEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		likeService      likeService
		likeEventService likeEventService
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
				mockLikeService.EXPECT().
					Delete(ctx, mockTx, like.ID).
					Return(nil)
				mockLikeEventService.EXPECT().Deleted(ctx, mockTx, like.ID).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockLikeService.EXPECT().
					Delete(ctx, mockTx, like.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &LikeUseCase{
				likeService:      tt.fields.likeService,
				likeEventService: tt.fields.likeEventService,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLikeUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeService := NewMocklikeService(ctrl)
	mockLikeEventService := NewMocklikeEventService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	filter := entities.NewMockLikeFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listLikes := make([]entities.Like, 0, count)
	for i := uint64(0); i < count; i++ {
		listLikes = append(listLikes, entities.NewMockLike(t))
	}
	type fields struct {
		likeService      likeService
		likeEventService likeEventService
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
		filter entities.LikeFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Like
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeService.EXPECT().
					List(ctx, filter).
					Return(listLikes, count, nil)
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listLikes,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockLikeService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				likeService:      mockLikeService,
				likeEventService: mockLikeEventService,
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
			i := &LikeUseCase{
				likeService:      tt.fields.likeService,
				likeEventService: tt.fields.likeEventService,
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
