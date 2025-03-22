package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	mockEntities "github.com/mikalai-mitsin/example/internal/app/comment/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewCommentServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentUseCase := NewMockcommentUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		commentUseCase commentUseCase
		logger         logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.CommentServiceServer
	}{
		{
			name: "ok",
			args: args{
				commentUseCase: mockCommentUseCase,
				logger:         mockLogger,
			},
			want: &CommentServiceServer{
				commentUseCase: mockCommentUseCase,
				logger:         mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommentServiceServer(tt.args.commentUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentUseCase := NewMockcommentUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := mockEntities.NewCommentCreate(t)
	comment := mockEntities.NewComment(t)
	type fields struct {
		UnimplementedCommentServiceServer examplepb.UnimplementedCommentServiceServer
		commentUseCase                    commentUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.CommentCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCommentUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(comment, nil)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.CommentCreate{},
			},
			want:    decodeComment(comment),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockCommentUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(entities.Comment{}, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.CommentCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := CommentServiceServer{
				UnimplementedCommentServiceServer: tt.fields.UnimplementedCommentServiceServer,
				commentUseCase:                    tt.fields.commentUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentUseCase := NewMockcommentUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedCommentServiceServer examplepb.UnimplementedCommentServiceServer
		commentUseCase                    commentUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.CommentDelete
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
				mockCommentUseCase.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CommentDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockCommentUseCase.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CommentDelete{
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
			s := CommentServiceServer{
				UnimplementedCommentServiceServer: tt.fields.UnimplementedCommentServiceServer,
				commentUseCase:                    tt.fields.commentUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentUseCase := NewMockcommentUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	type fields struct {
		UnimplementedCommentServiceServer examplepb.UnimplementedCommentServiceServer
		commentUseCase                    commentUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.CommentGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCommentUseCase.EXPECT().Get(ctx, comment.ID).Return(comment, nil).Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CommentGet{
					Id: string(comment.ID),
				},
			},
			want:    decodeComment(comment),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockCommentUseCase.EXPECT().Get(ctx, comment.ID).
					Return(entities.Comment{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CommentGet{
					Id: string(comment.ID),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := CommentServiceServer{
				UnimplementedCommentServiceServer: tt.fields.UnimplementedCommentServiceServer,
				commentUseCase:                    tt.fields.commentUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentUseCase := NewMockcommentUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := mockEntities.NewCommentFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListComment{
		Items: make([]*examplepb.Comment, 0, int(count)),
		Count: count,
	}
	listComments := make([]entities.Comment, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mockEntities.NewComment(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listComments = append(listComments, a)
		response.Items = append(response.Items, decodeComment(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedCommentServiceServer examplepb.UnimplementedCommentServiceServer
		commentUseCase                    commentUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.CommentFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListComment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCommentUseCase.EXPECT().
					List(ctx, gomock.Any()).
					Return(listComments, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CommentFilter{
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
				mockCommentUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CommentFilter{
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
			s := CommentServiceServer{
				UnimplementedCommentServiceServer: tt.fields.UnimplementedCommentServiceServer,
				commentUseCase:                    tt.fields.commentUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentUseCase := NewMockcommentUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	update := mockEntities.NewCommentUpdate(t)
	type fields struct {
		UnimplementedCommentServiceServer examplepb.UnimplementedCommentServiceServer
		commentUseCase                    commentUseCase
		logger                            logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.CommentUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCommentUseCase.EXPECT().Update(ctx, gomock.Any()).Return(comment, nil).Times(1)
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeCommentUpdate(update),
			},
			want:    decodeComment(comment),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockCommentUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(entities.Comment{}, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedCommentServiceServer: examplepb.UnimplementedCommentServiceServer{},
				commentUseCase:                    mockCommentUseCase,
				logger:                            mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeCommentUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := CommentServiceServer{
				UnimplementedCommentServiceServer: tt.fields.UnimplementedCommentServiceServer,
				commentUseCase:                    tt.fields.commentUseCase,
				logger:                            tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeComment(t *testing.T) {
	comment := mockEntities.NewComment(t)
	result := &examplepb.Comment{
		Id:        string(comment.ID),
		UpdatedAt: timestamppb.New(comment.UpdatedAt),
		CreatedAt: timestamppb.New(comment.CreatedAt),
		Text:      string(comment.Text),
		AuthorId:  string(comment.AuthorId),
		PostId:    string(comment.PostId),
	}
	type args struct {
		comment entities.Comment
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Comment
	}{
		{
			name: "ok",
			args: args{
				comment: comment,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeComment(tt.args.comment)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeCommentFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.CommentFilter
	}
	tests := []struct {
		name string
		args args
		want entities.CommentFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.CommentFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: entities.CommentFilter{
				PageSize:   pointer.Pointer(uint64(5)),
				PageNumber: pointer.Pointer(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				IDs:        []uuid.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeCommentFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
