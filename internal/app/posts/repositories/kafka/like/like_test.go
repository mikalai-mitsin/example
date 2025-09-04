package events

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/IBM/sarama"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewLikeEventProducer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	type args struct {
		producer producer
		logger   logger
	}
	tests := []struct {
		name string
		args args
		want *LikeEventProducer
	}{
		{
			name: "ok",
			args: args{
				producer: mockProducer,
				logger:   mockLogger,
			},
			want: &LikeEventProducer{
				producer: mockProducer,
				logger:   mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLikeEventProducer(tt.args.producer, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewLikeEventProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLikeEventProducer_Created(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx  context.Context
		dtx  dtx.TX
		like entities.Like
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				producer: mockProducer,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				like: like,
			},
			setup: func() {
				data, _ := json.Marshal(like)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventCreated,
					Value: data,
					Key:   like.ID.String(),
				}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "send error",
			fields: fields{
				producer: mockProducer,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				like: like,
			},
			setup: func() {
				data, _ := json.Marshal(like)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventCreated,
					Value: data,
					Key:   like.ID.String(),
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &LikeEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Created(tt.args.ctx, tt.args.dtx, tt.args.like)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLikeEventProducer_Updated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx  context.Context
		dtx  dtx.TX
		like entities.Like
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				producer: mockProducer,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				like: like,
			},
			setup: func() {
				data, _ := json.Marshal(like)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventUpdated,
					Value: data,
					Key:   like.ID.String(),
				}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "send error",
			fields: fields{
				producer: mockProducer,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				like: like,
			},
			setup: func() {
				data, _ := json.Marshal(like)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventUpdated,
					Value: data,
					Key:   like.ID.String(),
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &LikeEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Updated(tt.args.ctx, tt.args.dtx, tt.args.like)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLikeEventProducer_Deleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	like := entities.NewMockLike(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx context.Context
		dtx dtx.TX
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				producer: mockProducer,
				logger:   mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			setup: func() {
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventDeleted,
					Value: sarama.ByteEncoder(like.ID.String()),
					Key:   like.ID.String(),
				}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "send error",
			fields: fields{
				producer: mockProducer,
				logger:   mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			setup: func() {
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventDeleted,
					Value: []byte(like.ID.String()),
					Key:   like.ID.String(),
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &LikeEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Deleted(tt.args.ctx, tt.args.dtx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
