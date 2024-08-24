package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_interceptors "github.com/mikalai-mitsin/example/internal/app/arch/interceptors/mock"
	"github.com/mikalai-mitsin/example/internal/app/arch/models"
	mock_models "github.com/mikalai-mitsin/example/internal/app/arch/models/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewArchInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archUseCase := mock_interceptors.NewMockArchUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	type args struct {
		archUseCase ArchUseCase
		logger      Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *ArchInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				archUseCase: archUseCase,
				logger:      logger,
			},
			want: &ArchInterceptor{
				archUseCase: archUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewArchInterceptor(tt.args.archUseCase, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewArchInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archUseCase := mock_interceptors.NewMockArchUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	type fields struct {
		archUseCase ArchUseCase
		logger      Logger
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
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archUseCase.EXPECT().
					Get(ctx, arch.ID).
					Return(arch, nil)
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(arch.ID),
			},
			want:    arch,
			wantErr: nil,
		},
		{
			name: "Arch not found",
			setup: func() {
				archUseCase.EXPECT().
					Get(ctx, arch.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ArchInterceptor{
				archUseCase: tt.fields.archUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archUseCase := mock_interceptors.NewMockArchUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	create := mock_models.NewArchCreate(t)
	type fields struct {
		archUseCase ArchUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		create *models.ArchCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archUseCase.EXPECT().Create(ctx, create).Return(arch, nil)
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    arch,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				archUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				archUseCase: archUseCase,
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
			i := &ArchInterceptor{
				archUseCase: tt.fields.archUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archUseCase := mock_interceptors.NewMockArchUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	update := mock_models.NewArchUpdate(t)
	type fields struct {
		archUseCase ArchUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		update *models.ArchUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archUseCase.EXPECT().Update(ctx, update).Return(arch, nil)
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    arch,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				archUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
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
			i := &ArchInterceptor{
				archUseCase: tt.fields.archUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archUseCase := mock_interceptors.NewMockArchUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	type fields struct {
		archUseCase ArchUseCase
		logger      Logger
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
				archUseCase.EXPECT().
					Delete(ctx, arch.ID).
					Return(nil)
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				archUseCase.EXPECT().
					Delete(ctx, arch.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &ArchInterceptor{
				archUseCase: tt.fields.archUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestArchInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archUseCase := mock_interceptors.NewMockArchUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewArchFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listArches := make([]*models.Arch, 0, count)
	for i := uint64(0); i < count; i++ {
		listArches = append(listArches, mock_models.NewArch(t))
	}
	type fields struct {
		archUseCase ArchUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ArchFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Arch
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archUseCase.EXPECT().
					List(ctx, filter).
					Return(listArches, count, nil)
			},
			fields: fields{
				archUseCase: archUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listArches,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				archUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				archUseCase: archUseCase,
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
			i := &ArchInterceptor{
				archUseCase: tt.fields.archUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ArchInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
