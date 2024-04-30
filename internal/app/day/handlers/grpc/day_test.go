package grpc

import (
	"context"
	"errors"

	mock_grpc "github.com/018bf/example/internal/app/day/handlers/grpc/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"

	"reflect"
	"testing"

	"github.com/018bf/example/internal/app/day/models"
	mock_models "github.com/018bf/example/internal/app/day/models/mock"
	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewDayServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_grpc.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		dayInterceptor DayInterceptor
		logger         log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.DayServiceServer
	}{
		{
			name: "ok",
			args: args{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
			want: &DayServiceServer{
				dayInterceptor: dayInterceptor,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDayServiceServer(tt.args.dayInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewDayServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_grpc.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	// create := mock_models.NewDayCreate(t)
	day := mock_models.NewDay(t)
	type fields struct {
		UnimplementedDayServiceServer examplepb.UnimplementedDayServiceServer
		dayInterceptor                DayInterceptor
		logger                        log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.DayCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				dayInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(day, nil)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.DayCreate{},
			},
			want:    decodeDay(day),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				dayInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.DayCreate{},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := DayServiceServer{
				UnimplementedDayServiceServer: tt.fields.UnimplementedDayServiceServer,
				dayInterceptor:                tt.fields.dayInterceptor,
				logger:                        tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_grpc.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedDayServiceServer examplepb.UnimplementedDayServiceServer
		dayInterceptor                DayInterceptor
		logger                        log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.DayDelete
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
				dayInterceptor.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.DayDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				dayInterceptor.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.DayDelete{
					Id: id.String(),
				},
			},
			want: nil,
			wantErr: grpc.DecodeError(&errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "details", Value: "i error"}},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := DayServiceServer{
				UnimplementedDayServiceServer: tt.fields.UnimplementedDayServiceServer,
				dayInterceptor:                tt.fields.dayInterceptor,
				logger:                        tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_grpc.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		UnimplementedDayServiceServer examplepb.UnimplementedDayServiceServer
		dayInterceptor                DayInterceptor
		logger                        log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.DayGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				dayInterceptor.EXPECT().Get(ctx, day.ID).Return(day, nil).Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.DayGet{
					Id: string(day.ID),
				},
			},
			want:    decodeDay(day),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				dayInterceptor.EXPECT().Get(ctx, day.ID).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.DayGet{
					Id: string(day.ID),
				},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := DayServiceServer{
				UnimplementedDayServiceServer: tt.fields.UnimplementedDayServiceServer,
				dayInterceptor:                tt.fields.dayInterceptor,
				logger:                        tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_grpc.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewDayFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListDay{
		Items: make([]*examplepb.Day, 0, int(count)),
		Count: count,
	}
	listDays := make([]*models.Day, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewDay(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listDays = append(listDays, a)
		response.Items = append(response.Items, decodeDay(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedDayServiceServer examplepb.UnimplementedDayServiceServer
		dayInterceptor                DayInterceptor
		logger                        log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.DayFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListDay
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				dayInterceptor.EXPECT().List(ctx, filter).Return(listDays, count, nil).Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.DayFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				dayInterceptor.
					EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.DayFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := DayServiceServer{
				UnimplementedDayServiceServer: tt.fields.UnimplementedDayServiceServer,
				dayInterceptor:                tt.fields.dayInterceptor,
				logger:                        tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dayInterceptor := mock_grpc.NewMockDayInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	update := mock_models.NewDayUpdate(t)
	type fields struct {
		UnimplementedDayServiceServer examplepb.UnimplementedDayServiceServer
		dayInterceptor                DayInterceptor
		logger                        log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.DayUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				dayInterceptor.EXPECT().Update(ctx, gomock.Any()).Return(day, nil).Times(1)
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeDayUpdate(update),
			},
			want:    decodeDay(day),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				dayInterceptor.EXPECT().Update(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedDayServiceServer: examplepb.UnimplementedDayServiceServer{},
				dayInterceptor:                dayInterceptor,
				logger:                        logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeDayUpdate(update),
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := DayServiceServer{
				UnimplementedDayServiceServer: tt.fields.UnimplementedDayServiceServer,
				dayInterceptor:                tt.fields.dayInterceptor,
				logger:                        tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeDay(t *testing.T) {
	day := mock_models.NewDay(t)
	result := &examplepb.Day{
		Id:          string(day.ID),
		UpdatedAt:   timestamppb.New(day.UpdatedAt),
		CreatedAt:   timestamppb.New(day.CreatedAt),
		Name:        string(day.Name),
		Repeat:      int32(day.Repeat),
		EquipmentId: string(day.EquipmentID),
	}
	type args struct {
		day *models.Day
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Day
	}{
		{
			name: "ok",
			args: args{
				day: day,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeDay(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeDayFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.DayFilter
	}
	tests := []struct {
		name string
		args args
		want *models.DayFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.DayFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.DayFilter{
				PageSize:   pointer.Pointer(uint64(5)),
				PageNumber: pointer.Pointer(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				Search:     pointer.Pointer("my name is"),
				IDs:        []uuid.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeDayFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
