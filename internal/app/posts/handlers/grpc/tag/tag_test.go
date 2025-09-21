package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewTagServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagUseCase := NewMocktagUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		tagUseCase tagUseCase
		logger     logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.TagServiceServer
	}{
		{
			name: "ok",
			args: args{
				tagUseCase: mockTagUseCase,
				logger:     mockLogger,
			},
			want: &TagServiceServer{
				tagUseCase: mockTagUseCase,
				logger:     mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTagServiceServer(tt.args.tagUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagUseCase := NewMocktagUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := entities.NewMockTagCreate(t)
	tag := entities.NewMockTag(t)
	type fields struct {
		UnimplementedTagServiceServer examplepb.UnimplementedTagServiceServer
		tagUseCase                    tagUseCase
		logger                        logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.TagCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(tag, nil)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.TagCreate{},
			},
			want:    decodeTag(tag),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockTagUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(entities.Tag{}, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.TagCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := TagServiceServer{
				UnimplementedTagServiceServer: tt.fields.UnimplementedTagServiceServer,
				tagUseCase:                    tt.fields.tagUseCase,
				logger:                        tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagUseCase := NewMocktagUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	del := entities.NewMockTagDelete(t)
	del.ID = tag.ID
	type fields struct {
		UnimplementedTagServiceServer examplepb.UnimplementedTagServiceServer
		tagUseCase                    tagUseCase
		logger                        logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.TagDelete
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagUseCase.EXPECT().Delete(ctx, del).Return(tag, nil)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.TagDelete{
					Id: tag.ID.String(),
				},
			},
			want:    decodeTag(tag),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockTagUseCase.EXPECT().Delete(ctx, del).
					Return(entities.Tag{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.TagDelete{
					Id: tag.ID.String(),
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
			s := TagServiceServer{
				UnimplementedTagServiceServer: tt.fields.UnimplementedTagServiceServer,
				tagUseCase:                    tt.fields.tagUseCase,
				logger:                        tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagUseCase := NewMocktagUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	type fields struct {
		UnimplementedTagServiceServer examplepb.UnimplementedTagServiceServer
		tagUseCase                    tagUseCase
		logger                        logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.TagGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagUseCase.EXPECT().Get(ctx, tag.ID).Return(tag, nil).Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.TagGet{
					Id: tag.ID.String(),
				},
			},
			want:    decodeTag(tag),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockTagUseCase.EXPECT().Get(ctx, tag.ID).
					Return(entities.Tag{}, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.TagGet{
					Id: tag.ID.String(),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := TagServiceServer{
				UnimplementedTagServiceServer: tt.fields.UnimplementedTagServiceServer,
				tagUseCase:                    tt.fields.tagUseCase,
				logger:                        tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagUseCase := NewMocktagUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := entities.NewMockTagFilter(t)
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListTag{
		Items: make([]*examplepb.Tag, 0, int(count)),
		Count: count,
	}
	tags := make([]entities.Tag, 0, int(count))
	type fields struct {
		UnimplementedTagServiceServer examplepb.UnimplementedTagServiceServer
		tagUseCase                    tagUseCase
		logger                        logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.TagFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListTag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagUseCase.EXPECT().List(ctx, gomock.Any()).Return(tags, count, nil).Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.TagFilter{
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
				mockTagUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.TagFilter{
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
			s := TagServiceServer{
				UnimplementedTagServiceServer: tt.fields.UnimplementedTagServiceServer,
				tagUseCase:                    tt.fields.tagUseCase,
				logger:                        tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagUseCase := NewMocktagUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	update := entities.NewMockTagUpdate(t)
	type fields struct {
		UnimplementedTagServiceServer examplepb.UnimplementedTagServiceServer
		tagUseCase                    tagUseCase
		logger                        logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.TagUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagUseCase.EXPECT().Update(ctx, gomock.Any()).Return(tag, nil).Times(1)
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeTagUpdate(update),
			},
			want:    decodeTag(tag),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockTagUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(entities.Tag{}, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedTagServiceServer: examplepb.UnimplementedTagServiceServer{},
				tagUseCase:                    mockTagUseCase,
				logger:                        mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeTagUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := TagServiceServer{
				UnimplementedTagServiceServer: tt.fields.UnimplementedTagServiceServer,
				tagUseCase:                    tt.fields.tagUseCase,
				logger:                        tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeTag(t *testing.T) {
	tag := entities.NewMockTag(t)
	result := &examplepb.Tag{
		Id:        tag.ID.String(),
		UpdatedAt: timestamppb.New(tag.UpdatedAt),
		CreatedAt: timestamppb.New(tag.CreatedAt),
		DeletedAt: timestamppb.New(*tag.DeletedAt),
		PostId:    tag.PostId.String(),
		Value:     string(tag.Value),
	}
	type args struct {
		tag entities.Tag
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Tag
	}{
		{
			name: "ok",
			args: args{
				tag: tag,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeTag(tt.args.tag)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeTagFilter(t *testing.T) {
	type args struct {
		input *examplepb.TagFilter
	}
	tests := []struct {
		name string
		args args
		want entities.TagFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.TagFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					OrderBy:    []string{"created_at", "id"},
				},
			},
			want: entities.TagFilter{
				PageSize:   pointer.Of(uint64(5)),
				PageNumber: pointer.Of(uint64(2)),
				OrderBy:    []entities.TagOrdering{"created_at", "id"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeTagFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
