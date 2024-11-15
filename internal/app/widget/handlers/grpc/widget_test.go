package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/widget/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewWidgetServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetUseCase := NewMockwidgetUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		widgetUseCase widgetUseCase
		logger        logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.WidgetServiceServer
	}{
		{
			name: "ok",
			args: args{
				widgetUseCase: mockWidgetUseCase,
				logger:        mockLogger,
			},
			want: &WidgetServiceServer{
				widgetUseCase: mockWidgetUseCase,
				logger:        mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWidgetServiceServer(tt.args.widgetUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetUseCase := NewMockwidgetUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := mock_entities.NewWidgetCreate(t)
	widget := mock_entities.NewWidget(t)
	type fields struct {
		UnimplementedWidgetServiceServer examplepb.UnimplementedWidgetServiceServer
		widgetUseCase                    widgetUseCase
		logger                           logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.WidgetCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(widget, nil)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.WidgetCreate{},
			},
			want:    decodeWidget(widget),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockWidgetUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.WidgetCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := WidgetServiceServer{
				UnimplementedWidgetServiceServer: tt.fields.UnimplementedWidgetServiceServer,
				widgetUseCase:                    tt.fields.widgetUseCase,
				logger:                           tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetUseCase := NewMockwidgetUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedWidgetServiceServer examplepb.UnimplementedWidgetServiceServer
		widgetUseCase                    widgetUseCase
		logger                           logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.WidgetDelete
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
				mockWidgetUseCase.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.WidgetDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockWidgetUseCase.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.WidgetDelete{
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
			s := WidgetServiceServer{
				UnimplementedWidgetServiceServer: tt.fields.UnimplementedWidgetServiceServer,
				widgetUseCase:                    tt.fields.widgetUseCase,
				logger:                           tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetUseCase := NewMockwidgetUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	type fields struct {
		UnimplementedWidgetServiceServer examplepb.UnimplementedWidgetServiceServer
		widgetUseCase                    widgetUseCase
		logger                           logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.WidgetGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetUseCase.EXPECT().Get(ctx, widget.ID).Return(widget, nil).Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.WidgetGet{
					Id: string(widget.ID),
				},
			},
			want:    decodeWidget(widget),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockWidgetUseCase.EXPECT().Get(ctx, widget.ID).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.WidgetGet{
					Id: string(widget.ID),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := WidgetServiceServer{
				UnimplementedWidgetServiceServer: tt.fields.UnimplementedWidgetServiceServer,
				widgetUseCase:                    tt.fields.widgetUseCase,
				logger:                           tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetUseCase := NewMockwidgetUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := mock_entities.NewWidgetFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListWidget{
		Items: make([]*examplepb.Widget, 0, int(count)),
		Count: count,
	}
	listWidgets := make([]*entities.Widget, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_entities.NewWidget(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listWidgets = append(listWidgets, a)
		response.Items = append(response.Items, decodeWidget(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedWidgetServiceServer examplepb.UnimplementedWidgetServiceServer
		widgetUseCase                    widgetUseCase
		logger                           logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.WidgetFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListWidget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetUseCase.EXPECT().
					List(ctx, gomock.Any()).
					Return(listWidgets, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.WidgetFilter{
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
				mockWidgetUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.WidgetFilter{
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
			s := WidgetServiceServer{
				UnimplementedWidgetServiceServer: tt.fields.UnimplementedWidgetServiceServer,
				widgetUseCase:                    tt.fields.widgetUseCase,
				logger:                           tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWidgetUseCase := NewMockwidgetUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	widget := mock_entities.NewWidget(t)
	update := mock_entities.NewWidgetUpdate(t)
	type fields struct {
		UnimplementedWidgetServiceServer examplepb.UnimplementedWidgetServiceServer
		widgetUseCase                    widgetUseCase
		logger                           logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.WidgetUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockWidgetUseCase.EXPECT().Update(ctx, gomock.Any()).Return(widget, nil).Times(1)
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeWidgetUpdate(update),
			},
			want:    decodeWidget(widget),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockWidgetUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedWidgetServiceServer: examplepb.UnimplementedWidgetServiceServer{},
				widgetUseCase:                    mockWidgetUseCase,
				logger:                           mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeWidgetUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := WidgetServiceServer{
				UnimplementedWidgetServiceServer: tt.fields.UnimplementedWidgetServiceServer,
				widgetUseCase:                    tt.fields.widgetUseCase,
				logger:                           tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeWidget(t *testing.T) {
	widget := mock_entities.NewWidget(t)
	result := &examplepb.Widget{
		Id:           string(widget.ID),
		UpdatedAt:    timestamppb.New(widget.UpdatedAt),
		CreatedAt:    timestamppb.New(widget.CreatedAt),
		FormScreenId: string(widget.FormScreenId),
		Name:         string(widget.Name),
		Ordering:     int64(widget.Ordering),
		IsOptional:   bool(widget.IsOptional),
		UiSettings:   string(widget.UiSettings),
		DeletedAt:    timestamppb.New(widget.DeletedAt),
	}
	type args struct {
		widget *entities.Widget
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Widget
	}{
		{
			name: "ok",
			args: args{
				widget: widget,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeWidget(tt.args.widget)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeWidgetFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.WidgetFilter
	}
	tests := []struct {
		name string
		args args
		want *entities.WidgetFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.WidgetFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &entities.WidgetFilter{
				PageSize:   pointer.Pointer(uint64(5)),
				PageNumber: pointer.Pointer(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				IDs:        []uuid.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeWidgetFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
