package usecases

import (
	"context"
	"testing"

	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/widget/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewWidgetUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetService := NewMockwidgetService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		widgetService widgetService
		logger        logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *WidgetUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			want: &WidgetUseCase{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewWidgetUseCase(tt.args.widgetService, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetService := NewMockwidgetService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	type fields struct {
		widgetService widgetService
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
		want    *entities.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetService.EXPECT().
					Get(ctx, widget.ID).
					Return(widget, nil)
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(widget.ID),
			},
			want:    widget,
			wantErr: nil,
		},
		{
			name: "Widget not found",
			setup: func() {
				mockWidgetService.EXPECT().
					Get(ctx, widget.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &WidgetUseCase{
				widgetService: tt.fields.widgetService,
				logger:        tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetService := NewMockwidgetService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	create := mock_entities.NewWidgetCreate(t)
	type fields struct {
		widgetService widgetService
		logger        logger
	}
	type args struct {
		ctx    context.Context
		create *entities.WidgetCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetService.EXPECT().Create(ctx, create).Return(widget, nil)
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    widget,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockWidgetService.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &WidgetUseCase{
				widgetService: tt.fields.widgetService,
				logger:        tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetService := NewMockwidgetService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	update := mock_entities.NewWidgetUpdate(t)
	type fields struct {
		widgetService widgetService
		logger        logger
	}
	type args struct {
		ctx    context.Context
		update *entities.WidgetUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetService.EXPECT().Update(ctx, update).Return(widget, nil)
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    widget,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockWidgetService.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &WidgetUseCase{
				widgetService: tt.fields.widgetService,
				logger:        tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetService := NewMockwidgetService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	type fields struct {
		widgetService widgetService
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
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetService.EXPECT().
					Delete(ctx, widget.ID).
					Return(nil)
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockWidgetService.EXPECT().
					Delete(ctx, widget.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &WidgetUseCase{
				widgetService: tt.fields.widgetService,
				logger:        tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestWidgetUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetService := NewMockwidgetService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := mock_entities.NewWidgetFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listWidgets := make([]*entities.Widget, 0, count)
	for i := uint64(0); i < count; i++ {
		listWidgets = append(listWidgets, mock_entities.NewWidget(t))
	}
	type fields struct {
		widgetService widgetService
		logger        logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.WidgetFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entities.Widget
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetService.EXPECT().
					List(ctx, filter).
					Return(listWidgets, count, nil)
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listWidgets,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockWidgetService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				widgetService: mockWidgetService,
				logger:        mockLogger,
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
			i := &WidgetUseCase{
				widgetService: tt.fields.widgetService,
				logger:        tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
