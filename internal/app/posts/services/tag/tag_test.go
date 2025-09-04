package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewTagService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagRepository := NewMocktagRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	type args struct {
		tagRepository tagRepository
		clock         clock
		logger        logger
		uuid          uuidGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *TagService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
				uuid:          mockUUID,
			},
			want: &TagService{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
				uuid:          mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewTagService(tt.args.tagRepository, tt.args.clock, tt.args.logger, tt.args.uuid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagRepository := NewMocktagRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	type fields struct {
		tagRepository tagRepository
		logger        logger
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
		want    entities.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagRepository.EXPECT().Get(ctx, tag.ID).Return(tag, nil)
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  tag.ID,
			},
			want:    tag,
			wantErr: nil,
		},
		{
			name: "Tag not found",
			setup: func() {
				mockTagRepository.EXPECT().
					Get(ctx, tag.ID).
					Return(entities.Tag{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  tag.ID,
			},
			want:    entities.Tag{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &TagService{
				tagRepository: tt.fields.tagRepository,
				logger:        tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagRepository := NewMocktagRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listTags []entities.Tag
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listTags = append(listTags, entities.NewMockTag(t))
	}
	filter := entities.NewMockTagFilter(t)
	type fields struct {
		tagRepository tagRepository
		logger        logger
	}
	type args struct {
		ctx    context.Context
		filter entities.TagFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Tag
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagRepository.EXPECT().List(ctx, filter).Return(listTags, nil)
				mockTagRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listTags,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockTagRepository.EXPECT().
					List(ctx, filter).
					Return([]entities.Tag{}, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
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
				mockTagRepository.EXPECT().List(ctx, filter).Return(listTags, nil)
				mockTagRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
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
			u := &TagService{
				tagRepository: tt.fields.tagRepository,
				logger:        tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestTagService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagRepository := NewMocktagRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	create := entities.NewMockTagCreate(t)
	now := time.Now().UTC()
	type fields struct {
		tagRepository tagRepository
		clock         clock
		logger        logger
		uuid          uuidGenerator
	}
	type args struct {
		ctx    context.Context
		tx     dtx.TX
		create entities.TagCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    entities.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().
					NewUUID().
					Return(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
				mockTagRepository.EXPECT().
					Create(
						ctx,
						mockTx,
						entities.Tag{
							ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
							PostId:    create.PostId,
							Value:     create.Value,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
				uuid:          mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: create,
			},
			want: entities.Tag{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				PostId:    create.PostId,
				Value:     create.Value,
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
				mockTagRepository.EXPECT().
					Create(
						ctx,
						mockTx,
						entities.Tag{
							ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
							PostId:    create.PostId,
							Value:     create.Value,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
				uuid:          mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: create,
			},
			want:    entities.Tag{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
				clock:         mockClock,
				uuid:          mockUUID,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				create: entities.TagCreate{},
			},
			want: entities.Tag{},
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "post_id", Value: "cannot be blank"},
				errs.Param{Key: "value", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &TagService{
				tagRepository: tt.fields.tagRepository,
				clock:         tt.fields.clock,
				logger:        tt.fields.logger,
				uuid:          tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.tx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagRepository := NewMocktagRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	mockClock := NewMockclock(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	update := entities.NewMockTagUpdate(t)
	now := time.Now().UTC()
	updatedTag := entities.Tag{
		ID:        tag.ID,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: now,

		PostId: *update.PostId,
		Value:  *update.Value,
	}
	type fields struct {
		tagRepository tagRepository
		clock         clock
		logger        logger
	}
	type args struct {
		ctx    context.Context
		tx     dtx.TX
		update entities.TagUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    entities.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockTagRepository.EXPECT().
					Get(ctx, update.ID).Return(tag, nil)
				mockTagRepository.EXPECT().
					Update(ctx, mockTx, updatedTag).Return(nil)
			},
			fields: fields{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    updatedTag,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockTagRepository.EXPECT().
					Get(ctx, update.ID).
					Return(tag, nil)
				mockTagRepository.EXPECT().
					Update(ctx, mockTx, updatedTag).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    entities.Tag{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Tag not found",
			setup: func() {
				mockTagRepository.EXPECT().
					Get(ctx, update.ID).
					Return(entities.Tag{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				tagRepository: mockTagRepository,
				clock:         mockClock,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				tx:     mockTx,
				update: update,
			},
			want:    entities.Tag{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &TagService{
				tagRepository: tt.fields.tagRepository,
				clock:         tt.fields.clock,
				logger:        tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.tx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagRepository := NewMocktagRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	type fields struct {
		tagRepository tagRepository
		logger        logger
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
				mockTagRepository.EXPECT().
					Delete(ctx, mockTx, tag.ID).
					Return(nil)
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				id:  tag.ID,
			},
			wantErr: nil,
		},
		{
			name: "Tag not found",
			setup: func() {
				mockTagRepository.EXPECT().
					Delete(ctx, mockTx, tag.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				tagRepository: mockTagRepository,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				tx:  mockTx,
				id:  tag.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &TagService{
				tagRepository: tt.fields.tagRepository,
				logger:        tt.fields.logger,
			}
			err := u.Delete(tt.args.ctx, tt.args.tx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
