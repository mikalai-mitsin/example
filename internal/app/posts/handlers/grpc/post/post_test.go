package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewPostServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUseCase := NewMockpostUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		postUseCase postUseCase
		logger      logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.PostServiceServer
	}{
		{
			name: "ok",
			args: args{
				postUseCase: mockPostUseCase,
				logger:      mockLogger,
			},
			want: &PostServiceServer{
				postUseCase: mockPostUseCase,
				logger:      mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPostServiceServer(tt.args.postUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUseCase := NewMockpostUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := entities.NewMockPostCreate(t)
	post := entities.NewMockPost(t)
	type fields struct {
		UnimplementedPostServiceServer examplepb.UnimplementedPostServiceServer
		postUseCase                    postUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PostCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(post, nil)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.PostCreate{},
			},
			want:    decodePost(post),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockPostUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.PostCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PostServiceServer{
				UnimplementedPostServiceServer: tt.fields.UnimplementedPostServiceServer,
				postUseCase:                    tt.fields.postUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUseCase := NewMockpostUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedPostServiceServer examplepb.UnimplementedPostServiceServer
		postUseCase                    postUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PostDelete
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
				mockPostUseCase.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PostDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockPostUseCase.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PostDelete{
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
			s := PostServiceServer{
				UnimplementedPostServiceServer: tt.fields.UnimplementedPostServiceServer,
				postUseCase:                    tt.fields.postUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUseCase := NewMockpostUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		UnimplementedPostServiceServer examplepb.UnimplementedPostServiceServer
		postUseCase                    postUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PostGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostUseCase.EXPECT().Get(ctx, post.ID).Return(post, nil).Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PostGet{
					Id: post.ID.String(),
				},
			},
			want:    decodePost(post),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockPostUseCase.EXPECT().Get(ctx, post.ID).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PostGet{
					Id: post.ID.String(),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := PostServiceServer{
				UnimplementedPostServiceServer: tt.fields.UnimplementedPostServiceServer,
				postUseCase:                    tt.fields.postUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUseCase := NewMockpostUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := entities.NewMockPostFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListPost{
		Items: make([]*examplepb.Post, 0, int(count)),
		Count: count,
	}
	listPosts := make([]entities.Post, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := entities.NewMockPost(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, a.ID.String())
		listPosts = append(listPosts, a)
		response.Items = append(response.Items, decodePost(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedPostServiceServer examplepb.UnimplementedPostServiceServer
		postUseCase                    postUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PostFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListPost
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostUseCase.EXPECT().
					List(ctx, gomock.Any()).
					Return(listPosts, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PostFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockPostUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PostFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PostServiceServer{
				UnimplementedPostServiceServer: tt.fields.UnimplementedPostServiceServer,
				postUseCase:                    tt.fields.postUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUseCase := NewMockpostUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	update := entities.NewMockPostUpdate(t)
	type fields struct {
		UnimplementedPostServiceServer examplepb.UnimplementedPostServiceServer
		postUseCase                    postUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PostUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostUseCase.EXPECT().Update(ctx, gomock.Any()).Return(post, nil).Times(1)
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodePostUpdate(update),
			},
			want:    decodePost(post),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockPostUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedPostServiceServer: examplepb.UnimplementedPostServiceServer{},
				postUseCase:                    mockPostUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodePostUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PostServiceServer{
				UnimplementedPostServiceServer: tt.fields.UnimplementedPostServiceServer,
				postUseCase:                    tt.fields.postUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodePost(t *testing.T) {
	post := entities.NewMockPost(t)
	result := &examplepb.Post{
		Id:        post.ID.String(),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
		CreatedAt: timestamppb.New(post.CreatedAt),
		Body:      string(post.Body),
	}
	type args struct {
		post entities.Post
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Post
	}{
		{
			name: "ok",
			args: args{
				post: post,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodePost(tt.args.post)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodePostFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.PostFilter
	}
	tests := []struct {
		name string
		args args
		want entities.PostFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.PostFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{id.String()},
				},
			},
			want: entities.PostFilter{
				PageSize:   pointer.Of(uint64(5)),
				PageNumber: pointer.Of(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				IDs:        []uuid.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodePostFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
