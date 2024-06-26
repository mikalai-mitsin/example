package grpc

import (
	"context"
	"errors"

	mock_grpc "github.com/018bf/example/internal/app/plan/handlers/grpc/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"

	"reflect"
	"testing"

	"github.com/018bf/example/internal/app/plan/models"
	mock_models "github.com/018bf/example/internal/app/plan/models/mock"
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

func TestNewPlanServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_grpc.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		planInterceptor PlanInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.PlanServiceServer
	}{
		{
			name: "ok",
			args: args{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
			want: &PlanServiceServer{
				planInterceptor: planInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlanServiceServer(tt.args.planInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPlanServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_grpc.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	// create := mock_models.NewPlanCreate(t)
	plan := mock_models.NewPlan(t)
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(plan, nil)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.PlanCreate{},
			},
			want:    decodePlan(plan),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.PlanCreate{},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
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

func TestPlanServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_grpc.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanDelete
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
				planInterceptor.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanDelete{
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
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
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

func TestPlanServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_grpc.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.EXPECT().Get(ctx, plan.ID).Return(plan, nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanGet{
					Id: string(plan.ID),
				},
			},
			want:    decodePlan(plan),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.EXPECT().Get(ctx, plan.ID).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanGet{
					Id: string(plan.ID),
				},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
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

func TestPlanServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_grpc.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewPlanFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListPlan{
		Items: make([]*examplepb.Plan, 0, int(count)),
		Count: count,
	}
	listPlans := make([]*models.Plan, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewPlan(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listPlans = append(listPlans, a)
		response.Items = append(response.Items, decodePlan(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListPlan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.EXPECT().List(ctx, filter).Return(listPlans, count, nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanFilter{
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
				planInterceptor.
					EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.PlanFilter{
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
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
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

func TestPlanServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	planInterceptor := mock_grpc.NewMockPlanInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	plan := mock_models.NewPlan(t)
	update := mock_models.NewPlanUpdate(t)
	type fields struct {
		UnimplementedPlanServiceServer examplepb.UnimplementedPlanServiceServer
		planInterceptor                PlanInterceptor
		logger                         log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.PlanUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Plan
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				planInterceptor.EXPECT().Update(ctx, gomock.Any()).Return(plan, nil).Times(1)
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodePlanUpdate(update),
			},
			want:    decodePlan(plan),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				planInterceptor.EXPECT().Update(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedPlanServiceServer: examplepb.UnimplementedPlanServiceServer{},
				planInterceptor:                planInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodePlanUpdate(update),
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := PlanServiceServer{
				UnimplementedPlanServiceServer: tt.fields.UnimplementedPlanServiceServer,
				planInterceptor:                tt.fields.planInterceptor,
				logger:                         tt.fields.logger,
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

func Test_decodePlan(t *testing.T) {
	plan := mock_models.NewPlan(t)
	result := &examplepb.Plan{
		Id:          string(plan.ID),
		UpdatedAt:   timestamppb.New(plan.UpdatedAt),
		CreatedAt:   timestamppb.New(plan.CreatedAt),
		Name:        string(plan.Name),
		Repeat:      uint64(plan.Repeat),
		EquipmentId: string(plan.EquipmentID),
	}
	type args struct {
		plan *models.Plan
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Plan
	}{
		{
			name: "ok",
			args: args{
				plan: plan,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodePlan(tt.args.plan); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodePlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodePlanFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.PlanFilter
	}
	tests := []struct {
		name string
		args args
		want *models.PlanFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.PlanFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.PlanFilter{
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
			if got := encodePlanFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
