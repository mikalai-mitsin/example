package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/widget/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewWidgetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetRepository := NewMockwidgetRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockUUIDGenerator(ctrl)
	type args struct {
		widgetRepository widgetRepository
		clock            clock
		logger           logger
		uuid             UUIDGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *WidgetService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
				uuid:             mockUUID,
			},
			want: &WidgetService{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
				uuid:             mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewWidgetService(
				tt.args.widgetRepository,
				tt.args.clock,
				tt.args.logger,
				tt.args.uuid,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetRepository := NewMockwidgetRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	type fields struct {
		widgetRepository widgetRepository
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
		want    *entities.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetRepository.EXPECT().Get(ctx, widget.ID).Return(widget, nil)
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			want:    widget,
			wantErr: nil,
		},
		{
			name: "Widget not found",
			setup: func() {
				mockWidgetRepository.EXPECT().
					Get(ctx, widget.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
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
			u := &WidgetService{
				widgetRepository: tt.fields.widgetRepository,
				logger:           tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetRepository := NewMockwidgetRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listWidgets []*entities.Widget
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listWidgets = append(listWidgets, mock_entities.NewWidget(t))
	}
	filter := mock_entities.NewWidgetFilter(t)
	type fields struct {
		widgetRepository widgetRepository
		logger           logger
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
				mockWidgetRepository.EXPECT().List(ctx, filter).Return(listWidgets, nil)
				mockWidgetRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
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
				mockWidgetRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
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
				mockWidgetRepository.EXPECT().List(ctx, filter).Return(listWidgets, nil)
				mockWidgetRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
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
			u := &WidgetService{
				widgetRepository: tt.fields.widgetRepository,
				logger:           tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestWidgetService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetRepository := NewMockwidgetRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockUUIDGenerator(ctrl)
	ctx := context.Background()
	create := mock_entities.NewWidgetCreate(t)
	now := time.Now().UTC()
	type fields struct {
		widgetRepository widgetRepository
		clock            clock
		logger           logger
		uuid             UUIDGenerator
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
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test"))
				mockWidgetRepository.EXPECT().
					Create(
						ctx,
						&entities.Widget{
							ID:           uuid.UUID("test"),
							FormScreenId: create.FormScreenId,
							Name:         create.Name,
							Ordering:     create.Ordering,
							IsOptional:   create.IsOptional,
							UiSettings:   create.UiSettings,
							DeletedAt:    create.DeletedAt,
							UpdatedAt:    now,
							CreatedAt:    now,
						},
					).
					Return(nil)
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
				uuid:             mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &entities.Widget{
				ID:           uuid.UUID("test"),
				FormScreenId: create.FormScreenId,
				Name:         create.Name,
				Ordering:     create.Ordering,
				IsOptional:   create.IsOptional,
				UiSettings:   create.UiSettings,
				DeletedAt:    create.DeletedAt,
				UpdatedAt:    now,
				CreatedAt:    now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test 2"))
				mockWidgetRepository.EXPECT().
					Create(
						ctx,
						&entities.Widget{
							ID:           uuid.UUID("test 2"),
							FormScreenId: create.FormScreenId,
							Name:         create.Name,
							Ordering:     create.Ordering,
							IsOptional:   create.IsOptional,
							UiSettings:   create.UiSettings,
							DeletedAt:    create.DeletedAt,
							UpdatedAt:    now,
							CreatedAt:    now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
				uuid:             mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
				clock:            mockClock,
				uuid:             mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: &entities.WidgetCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "form_screen_id", Value: "cannot be blank"},
				errs.Param{Key: "name", Value: "cannot be blank"},
				errs.Param{Key: "ordering", Value: "cannot be blank"},
				errs.Param{Key: "is_optional", Value: "cannot be blank"},
				errs.Param{Key: "ui_settings", Value: "cannot be blank"},
				errs.Param{Key: "deleted_at", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &WidgetService{
				widgetRepository: tt.fields.widgetRepository,
				clock:            tt.fields.clock,
				logger:           tt.fields.logger,
				uuid:             tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetRepository := NewMockwidgetRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	mockClock := NewMockclock(ctrl)
	update := mock_entities.NewWidgetUpdate(t)
	now := widget.UpdatedAt
	type fields struct {
		widgetRepository widgetRepository
		clock            clock
		logger           logger
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
				mockClock.EXPECT().Now().Return(now)
				mockWidgetRepository.EXPECT().
					Get(ctx, update.ID).Return(widget, nil)
				mockWidgetRepository.EXPECT().
					Update(ctx, widget).Return(nil)
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
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
				mockClock.EXPECT().Now().Return(now)
				mockWidgetRepository.EXPECT().
					Get(ctx, update.ID).
					Return(widget, nil)
				mockWidgetRepository.EXPECT().
					Update(ctx, widget).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Widget not found",
			setup: func() {
				mockWidgetRepository.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				clock:            mockClock,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				update: &entities.WidgetUpdate{
					ID: uuid.UUID("baduuid"),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &WidgetService{
				widgetRepository: tt.fields.widgetRepository,
				clock:            tt.fields.clock,
				logger:           tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetRepository := NewMockwidgetRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	type fields struct {
		widgetRepository widgetRepository
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
				mockWidgetRepository.EXPECT().
					Delete(ctx, widget.ID).
					Return(nil)
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			wantErr: nil,
		},
		{
			name: "Widget not found",
			setup: func() {
				mockWidgetRepository.EXPECT().
					Delete(ctx, widget.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				widgetRepository: mockWidgetRepository,
				logger:           mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &WidgetService{
				widgetRepository: tt.fields.widgetRepository,
				logger:           tt.fields.logger,
			}
			err := u.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
