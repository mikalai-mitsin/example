package events

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/IBM/sarama"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewPostEventProducer(t *testing.T) {
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
		want *PostEventProducer
	}{
		{
			name: "ok",
			args: args{
				producer: mockProducer,
				logger:   mockLogger,
			},
			want: &PostEventProducer{
				producer: mockProducer,
				logger:   mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostEventProducer(tt.args.producer, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewPostEventProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostEventProducer_Created(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx  context.Context
		post entities.Post
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
				post: post,
			},
			setup: func() {
				data, _ := json.Marshal(post)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventCreated,
					Value: data,
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
				post: post,
			},
			setup: func() {
				data, _ := json.Marshal(post)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventCreated,
					Value: data,
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &PostEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Created(tt.args.ctx, tt.args.post)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostEventProducer_Updated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx  context.Context
		post entities.Post
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
				post: post,
			},
			setup: func() {
				data, _ := json.Marshal(post)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventUpdated,
					Value: data,
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
				post: post,
			},
			setup: func() {
				data, _ := json.Marshal(post)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventUpdated,
					Value: data,
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &PostEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Updated(tt.args.ctx, tt.args.post)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostEventProducer_Deleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx context.Context
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
				id:  post.ID,
			},
			setup: func() {
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventDeleted,
					Value: sarama.ByteEncoder(post.ID.String()),
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
				id:  post.ID,
			},
			setup: func() {
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventDeleted,
					Value: []byte(post.ID.String()),
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &PostEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Deleted(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
