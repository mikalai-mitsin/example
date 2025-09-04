package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewLikeService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeRepository := NewMocklikeRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	type args struct {
		likeRepository likeRepository
		clock          clock
		logger         logger
		uuid           uuidGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *LikeService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			want: &LikeService{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewLikeService(
				tt.args.likeRepository,
				tt.args.clock,
				tt.args.logger,
				tt.args.uuid,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeRepository := NewMocklikeRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		likeRepository likeRepository
		logger         logger
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
				mockLikeRepository.EXPECT().Get(ctx, like.ID).Return(like, nil)
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			want:    like,
			wantErr: nil,
		},
		{
			name: "Like not found",
			setup: func() {
				mockLikeRepository.EXPECT().
					Get(ctx, like.ID).
					Return(entities.Like{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
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
			u := &LikeService{
				likeRepository: tt.fields.likeRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeRepository := NewMocklikeRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listLikes []entities.Like
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listLikes = append(listLikes, entities.NewMockLike(t))
	}
	filter := entities.NewMockLikeFilter(t)
	type fields struct {
		likeRepository likeRepository
		logger         logger
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
				mockLikeRepository.EXPECT().List(ctx, filter).Return(listLikes, nil)
				mockLikeRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
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
				mockLikeRepository.EXPECT().
					List(ctx, filter).
					Return([]entities.Like{}, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
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
				mockLikeRepository.EXPECT().List(ctx, filter).Return(listLikes, nil)
				mockLikeRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
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
			u := &LikeService{
				likeRepository: tt.fields.likeRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestLikeService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeRepository := NewMocklikeRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	create := entities.NewMockLikeCreate(t)
	now := time.Now().UTC()
	type fields struct {
		likeRepository likeRepository
		clock          clock
		logger         logger
		uuid           uuidGenerator
	}
	type args struct {
		ctx    context.Context
		tx     dtx.TX
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
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().
					NewUUID().
					Return(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
				mockLikeRepository.EXPECT().
					Create(
						ctx,
						mockTx,
						entities.Like{
							ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							PostId:    create.PostId,
							Value:     create.Value,
							UserId:    create.UserId,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: create,
			},
			want: entities.Like{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostId:    create.PostId,
				Value:     create.Value,
				UserId:    create.UserId,
				UpdatedAt: now,
				CreatedAt: now,
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
				mockLikeRepository.EXPECT().
					Create(
						ctx,
						mockTx,
						entities.Like{
							ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							PostId:    create.PostId,
							Value:     create.Value,
							UserId:    create.UserId,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: create,
			},
			want:    entities.Like{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
				clock:          mockClock,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: entities.LikeCreate{},
			},
			want: entities.Like{},
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "post_id", Value: "cannot be blank"},
				errs.Param{Key: "value", Value: "cannot be blank"},
				errs.Param{Key: "user_id", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &LikeService{
				likeRepository: tt.fields.likeRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
				uuid:           tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.tx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeRepository := NewMocklikeRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	mockClock := NewMockclock(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	update := entities.NewMockLikeUpdate(t)
	now := time.Now().UTC()
	updatedLike := entities.Like{
		ID:        like.ID,
		CreatedAt: like.CreatedAt,
		UpdatedAt: now,

		PostId: *update.PostId,
		Value:  *update.Value,
		UserId: *update.UserId,
	}
	type fields struct {
		likeRepository likeRepository
		clock          clock
		logger         logger
	}
	type args struct {
		ctx    context.Context
		tx     dtx.TX
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
				mockClock.EXPECT().Now().Return(now)
				mockLikeRepository.EXPECT().
					Get(ctx, update.ID).Return(like, nil)
				mockLikeRepository.EXPECT().
					Update(ctx, mockTx, updatedLike).Return(nil)
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    updatedLike,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockLikeRepository.EXPECT().
					Get(ctx, update.ID).
					Return(like, nil)
				mockLikeRepository.EXPECT().
					Update(ctx, mockTx, updatedLike).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    entities.Like{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Like not found",
			setup: func() {
				mockLikeRepository.EXPECT().
					Get(ctx, update.ID).
					Return(entities.Like{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    entities.Like{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &LikeService{
				likeRepository: tt.fields.likeRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.tx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeRepository := NewMocklikeRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		likeRepository likeRepository
		logger         logger
	}
	type args struct {
		ctx context.Context
		tx  dtx.TX
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
				mockLikeRepository.EXPECT().
					Delete(ctx, mockTx, like.ID).
					Return(nil)
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				id:  like.ID,
			},
			wantErr: nil,
		},
		{
			name: "Like not found",
			setup: func() {
				mockLikeRepository.EXPECT().
					Delete(ctx, mockTx, like.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				likeRepository: mockLikeRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				id:  like.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &LikeService{
				likeRepository: tt.fields.likeRepository,
				logger:         tt.fields.logger,
			}
			err := u.Delete(tt.args.ctx, tt.args.tx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
