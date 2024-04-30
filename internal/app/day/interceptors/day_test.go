package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_interceptors "github.com/018bf/example/internal/app/day/interceptors/mock"
	"github.com/018bf/example/internal/app/day/models"
	mock_models "github.com/018bf/example/internal/app/day/models/mock"
	userModels "github.com/018bf/example/internal/app/user/models"
	userMockModels "github.com/018bf/example/internal/app/user/models/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"

	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"

	"github.com/018bf/example/internal/pkg/uuid"
)

func TestNewDayInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	dayUseCase := mock_interceptors.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase AuthUseCase
		dayUseCase  DayUseCase
		logger      log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *DayInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			want: &DayInterceptor{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewDayInterceptor(tt.args.dayUseCase, tt.args.logger, tt.args.authUseCase); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewDayInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	dayUseCase := mock_interceptors.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		authUseCase AuthUseCase
		dayUseCase  DayUseCase
		logger      log.Logger
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
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDetail).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDetail, day).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(day.ID),
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDetail).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDetail, day).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Day not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDetail).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	dayUseCase := mock_interceptors.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	create := mock_models.NewDayCreate(t)
	type fields struct {
		dayUseCase  DayUseCase
		authUseCase AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.DayCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayCreate, create).
					Return(nil)
				dayUseCase.EXPECT().Create(ctx, create).Return(day, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "create error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayCreate, create).
					Return(nil)
				dayUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
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
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	dayUseCase := mock_interceptors.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	update := mock_models.NewDayUpdate(t)
	type fields struct {
		dayUseCase  DayUseCase
		authUseCase AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.DayUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayUpdate, day).
					Return(nil)
				dayUseCase.EXPECT().Update(ctx, update).Return(day, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayUpdate, day).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayUpdate).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayUpdate, day).
					Return(nil)
				dayUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	dayUseCase := mock_interceptors.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	day := mock_models.NewDay(t)
	type fields struct {
		dayUseCase  DayUseCase
		authUseCase AuthUseCase
		logger      log.Logger
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
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDelete, day).
					Return(nil)
				dayUseCase.EXPECT().
					Delete(ctx, day.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: nil,
		},
		{
			name: "Day not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, errs.NewEntityNotFoundError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDelete, day).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDelete).
					Return(nil)
				dayUseCase.EXPECT().
					Get(ctx, day.ID).
					Return(day, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayDelete, day).
					Return(nil)
				dayUseCase.EXPECT().
					Delete(ctx, day.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				dayUseCase:  dayUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	dayUseCase := mock_interceptors.NewMockDayUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewDayFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listDays := make([]*models.Day, 0, count)
	for i := uint64(0); i < count; i++ {
		listDays = append(listDays, mock_models.NewDay(t))
	}
	type fields struct {
		dayUseCase  DayUseCase
		authUseCase AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.DayFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Day
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayList, filter).
					Return(nil)
				dayUseCase.EXPECT().
					List(ctx, filter).
					Return(listDays, count, nil)
			},
			fields: fields{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listDays,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "list error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDDayList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDDayList, filter).
					Return(nil)
				dayUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				dayUseCase:  dayUseCase,
				authUseCase: authUseCase,
				logger:      logger,
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
			i := &DayInterceptor{
				dayUseCase:  tt.fields.dayUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DayInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
