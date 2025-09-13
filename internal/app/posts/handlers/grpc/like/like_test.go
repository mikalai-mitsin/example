package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewLikeServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeUseCase := NewMocklikeUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		likeUseCase likeUseCase
		logger      logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.LikeServiceServer
	}{
		{
			name: "ok",
			args: args{
				likeUseCase: mockLikeUseCase,
				logger:      mockLogger,
			},
			want: &LikeServiceServer{
				likeUseCase: mockLikeUseCase,
				logger:      mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLikeServiceServer(tt.args.likeUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeUseCase := NewMocklikeUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := entities.NewMockLikeCreate(t)
	like := entities.NewMockLike(t)
	type fields struct {
		UnimplementedLikeServiceServer examplepb.UnimplementedLikeServiceServer
		likeUseCase                    likeUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.LikeCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(like, nil)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.LikeCreate{},
			},
			want:    decodeLike(like),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockLikeUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(entities.Like{}, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.LikeCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := LikeServiceServer{
				UnimplementedLikeServiceServer: tt.fields.UnimplementedLikeServiceServer,
				likeUseCase:                    tt.fields.likeUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeUseCase := NewMocklikeUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		UnimplementedLikeServiceServer examplepb.UnimplementedLikeServiceServer
		likeUseCase                    likeUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.LikeDelete
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeUseCase.EXPECT().Delete(ctx, like.ID).Return(like, nil)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.LikeDelete{
					Id: like.ID.String(),
				},
			},
			want:    decodeLike(like),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockLikeUseCase.EXPECT().Delete(ctx, like.ID).
					Return(entities.Like{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.LikeDelete{
					Id: like.ID.String(),
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
			s := LikeServiceServer{
				UnimplementedLikeServiceServer: tt.fields.UnimplementedLikeServiceServer,
				likeUseCase:                    tt.fields.likeUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeUseCase := NewMocklikeUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		UnimplementedLikeServiceServer examplepb.UnimplementedLikeServiceServer
		likeUseCase                    likeUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.LikeGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeUseCase.EXPECT().Get(ctx, like.ID).Return(like, nil).Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.LikeGet{
					Id: like.ID.String(),
				},
			},
			want:    decodeLike(like),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockLikeUseCase.EXPECT().Get(ctx, like.ID).
					Return(entities.Like{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.LikeGet{
					Id: like.ID.String(),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := LikeServiceServer{
				UnimplementedLikeServiceServer: tt.fields.UnimplementedLikeServiceServer,
				likeUseCase:                    tt.fields.likeUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeUseCase := NewMocklikeUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := entities.NewMockLikeFilter(t)
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListLike{
		Items: make([]*examplepb.Like, 0, int(count)),
		Count: count,
	}
	likes := make([]entities.Like, 0, int(count))
	type fields struct {
		UnimplementedLikeServiceServer examplepb.UnimplementedLikeServiceServer
		likeUseCase                    likeUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.LikeFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListLike
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeUseCase.EXPECT().List(ctx, gomock.Any()).Return(likes, count, nil).Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.LikeFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					OrderBy:    nil,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockLikeUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.LikeFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
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
			s := LikeServiceServer{
				UnimplementedLikeServiceServer: tt.fields.UnimplementedLikeServiceServer,
				likeUseCase:                    tt.fields.likeUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLikeUseCase := NewMocklikeUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	update := entities.NewMockLikeUpdate(t)
	type fields struct {
		UnimplementedLikeServiceServer examplepb.UnimplementedLikeServiceServer
		likeUseCase                    likeUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.LikeUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockLikeUseCase.EXPECT().Update(ctx, gomock.Any()).Return(like, nil).Times(1)
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeLikeUpdate(update),
			},
			want:    decodeLike(like),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockLikeUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(entities.Like{}, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedLikeServiceServer: examplepb.UnimplementedLikeServiceServer{},
				likeUseCase:                    mockLikeUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeLikeUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := LikeServiceServer{
				UnimplementedLikeServiceServer: tt.fields.UnimplementedLikeServiceServer,
				likeUseCase:                    tt.fields.likeUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeLike(t *testing.T) {
	like := entities.NewMockLike(t)
	result := &examplepb.Like{
		Id:        like.ID.String(),
		UpdatedAt: timestamppb.New(like.UpdatedAt),
		CreatedAt: timestamppb.New(like.CreatedAt),
		DeletedAt: timestamppb.New(*like.DeletedAt),
		PostId:    like.PostId.String(),
		Value:     string(like.Value),
		UserId:    like.UserId.String(),
	}
	type args struct {
		like entities.Like
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Like
	}{
		{
			name: "ok",
			args: args{
				like: like,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeLike(tt.args.like)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeLikeFilter(t *testing.T) {
	type args struct {
		input *examplepb.LikeFilter
	}
	tests := []struct {
		name string
		args args
		want entities.LikeFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.LikeFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					OrderBy:    []string{"created_at", "id"},
				},
			},
			want: entities.LikeFilter{
				PageSize:   pointer.Of(uint64(5)),
				PageNumber: pointer.Of(uint64(2)),
				OrderBy:    []entities.LikeOrdering{"created_at", "id"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeLikeFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
