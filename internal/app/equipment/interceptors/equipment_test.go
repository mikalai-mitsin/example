package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_interceptors "github.com/018bf/example/internal/app/equipment/interceptors/mock"
	"github.com/018bf/example/internal/app/equipment/models"
	mock_models "github.com/018bf/example/internal/app/equipment/models/mock"
	userModels "github.com/018bf/example/internal/app/user/models"
	userMockModels "github.com/018bf/example/internal/app/user/models/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"

	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"

	"github.com/018bf/example/internal/pkg/uuid"
)

func TestNewEquipmentInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	equipmentUseCase := mock_interceptors.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase      AuthUseCase
		equipmentUseCase EquipmentUseCase
		logger           log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *EquipmentInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			want: &EquipmentInterceptor{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewEquipmentInterceptor(tt.args.equipmentUseCase, tt.args.logger, tt.args.authUseCase); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewEquipmentInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	equipmentUseCase := mock_interceptors.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		authUseCase      AuthUseCase
		equipmentUseCase EquipmentUseCase
		logger           log.Logger
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
		want    *models.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail, equipment).
					Return(nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(equipment.ID),
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail, equipment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Equipment not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDetail).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	equipmentUseCase := mock_interceptors.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	create := mock_models.NewEquipmentCreate(t)
	type fields struct {
		equipmentUseCase EquipmentUseCase
		authUseCase      AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.EquipmentCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate, create).
					Return(nil)
				equipmentUseCase.EXPECT().Create(ctx, create).Return(equipment, nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentCreate, create).
					Return(nil)
				equipmentUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
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
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	equipmentUseCase := mock_interceptors.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	update := mock_models.NewEquipmentUpdate(t)
	type fields struct {
		equipmentUseCase EquipmentUseCase
		authUseCase      AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.EquipmentUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Equipment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().Update(ctx, update).Return(equipment, nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    equipment,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate, equipment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
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
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	equipmentUseCase := mock_interceptors.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	equipment := mock_models.NewEquipment(t)
	type fields struct {
		equipmentUseCase EquipmentUseCase
		authUseCase      AuthUseCase
		logger           log.Logger
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().
					Delete(ctx, equipment.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: nil,
		},
		{
			name: "Equipment not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete, equipment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete).
					Return(nil)
				equipmentUseCase.EXPECT().
					Get(ctx, equipment.ID).
					Return(equipment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete, equipment).
					Return(nil)
				equipmentUseCase.EXPECT().
					Delete(ctx, equipment.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:      authUseCase,
				equipmentUseCase: equipmentUseCase,
				logger:           logger,
			},
			args: args{
				ctx: ctx,
				id:  equipment.ID,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEquipmentInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	equipmentUseCase := mock_interceptors.NewMockEquipmentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewEquipmentFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listEquipment := make([]*models.Equipment, 0, count)
	for i := uint64(0); i < count; i++ {
		listEquipment = append(listEquipment, mock_models.NewEquipment(t))
	}
	type fields struct {
		equipmentUseCase EquipmentUseCase
		authUseCase      AuthUseCase
		logger           log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.EquipmentFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Equipment
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentList, filter).
					Return(nil)
				equipmentUseCase.EXPECT().
					List(ctx, filter).
					Return(listEquipment, count, nil)
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listEquipment,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
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
					HasPermission(ctx, requestUser, userModels.PermissionIDEquipmentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDEquipmentList, filter).
					Return(nil)
				equipmentUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				equipmentUseCase: equipmentUseCase,
				authUseCase:      authUseCase,
				logger:           logger,
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
			i := &EquipmentInterceptor{
				equipmentUseCase: tt.fields.equipmentUseCase,
				authUseCase:      tt.fields.authUseCase,
				logger:           tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EquipmentInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EquipmentInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("EquipmentInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
