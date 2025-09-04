package usecases

import (
	"context"
	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewTagUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagService := NewMocktagService(ctrl)
	mocktagEventProducer := NewMocktagEventProducer(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		tagService       tagService
		tagEventProducer tagEventProducer
		dtxManager       dtxManager
		logger           logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *TagUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			want: &TagUseCase{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewTagUseCase(
				tt.args.tagService,
				tt.args.tagEventProducer,
				tt.args.dtxManager,
				tt.args.logger,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagService := NewMocktagService(ctrl)
	mocktagEventProducer := NewMocktagEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	type fields struct {
		tagService       tagService
		tagEventProducer tagEventProducer
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
		want    entities.Tag
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockTagService.EXPECT().
					Get(ctx, tag.ID).
					Return(tag, nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(tag.ID),
			},
			want:    tag,
			wantErr: nil,
		},
		{
			name: "Tag not found",
			setup: func() {
				mockTagService.EXPECT().
					Get(ctx, tag.ID).
					Return(entities.Tag{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
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
			i := &TagUseCase{
				tagService:       tt.fields.tagService,
				tagEventProducer: tt.fields.tagEventProducer,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagService := NewMocktagService(ctrl)
	mocktagEventProducer := NewMocktagEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	create := entities.NewMockTagCreate(t)
	type fields struct {
		tagService       tagService
		tagEventProducer tagEventProducer
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockTagService.EXPECT().Create(ctx, mockTx, create).Return(tag, nil)
				mocktagEventProducer.EXPECT().Created(ctx, mockTx, tag).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    tag,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockTagService.EXPECT().
					Create(ctx, mockTx, create).
					Return(entities.Tag{}, errs.NewUnexpectedBehaviorError("c u"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Tag{},
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &TagUseCase{
				tagService:       tt.fields.tagService,
				tagEventProducer: tt.fields.tagEventProducer,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagService := NewMocktagService(ctrl)
	mocktagEventProducer := NewMocktagEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	update := entities.NewMockTagUpdate(t)
	type fields struct {
		tagService       tagService
		tagEventProducer tagEventProducer
		dtxManager       dtxManager
		logger           logger
	}
	type args struct {
		ctx    context.Context
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
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockTagService.EXPECT().Update(ctx, mockTx, update).Return(tag, nil)
				mocktagEventProducer.EXPECT().Updated(ctx, mockTx, tag).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    tag,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockTagService.EXPECT().
					Update(ctx, mockTx, update).
					Return(entities.Tag{}, errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Tag{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &TagUseCase{
				tagService:       tt.fields.tagService,
				tagEventProducer: tt.fields.tagEventProducer,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTagUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagService := NewMocktagService(ctrl)
	mocktagEventProducer := NewMocktagEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockLogger.EXPECT().WithContext(gomock.Any()).Return(mockLogger).AnyTimes()
	mockDtxManager := NewMockdtxManager(ctrl)
	mockTx := dtx.NewMockTX(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	type fields struct {
		tagService       tagService
		tagEventProducer tagEventProducer
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
				mockTagService.EXPECT().
					Delete(ctx, mockTx, tag.ID).
					Return(nil)
				mocktagEventProducer.EXPECT().Deleted(ctx, mockTx, tag.ID).Return(nil)
				mockTx.EXPECT().Rollback().After(mockTx.EXPECT().Commit().Return(nil)).Return(nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  tag.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockDtxManager.EXPECT().NewTx().Return(mockTx)
				mockTagService.EXPECT().
					Delete(ctx, mockTx, tag.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
				mockTx.EXPECT().Rollback().Return(nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  tag.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &TagUseCase{
				tagService:       tt.fields.tagService,
				tagEventProducer: tt.fields.tagEventProducer,
				dtxManager:       tt.fields.dtxManager,
				logger:           tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTagUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagService := NewMocktagService(ctrl)
	mocktagEventProducer := NewMocktagEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockDtxManager := NewMockdtxManager(ctrl)
	ctx := context.Background()
	filter := entities.NewMockTagFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listTags := make([]entities.Tag, 0, count)
	for i := uint64(0); i < count; i++ {
		listTags = append(listTags, entities.NewMockTag(t))
	}
	type fields struct {
		tagService       tagService
		tagEventProducer tagEventProducer
		dtxManager       dtxManager
		logger           logger
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
				mockTagService.EXPECT().
					List(ctx, filter).
					Return(listTags, count, nil)
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
				dtxManager:       mockDtxManager,
				logger:           mockLogger,
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
				mockTagService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				tagService:       mockTagService,
				tagEventProducer: mocktagEventProducer,
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
			i := &TagUseCase{
				tagService:       tt.fields.tagService,
				tagEventProducer: tt.fields.tagEventProducer,
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
