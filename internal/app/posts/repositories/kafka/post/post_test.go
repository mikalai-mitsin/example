package repositories

import (
	"context"
	"errors"
	"reflect"
	"testing"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
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

func TestPostEventProducer_Send(t *testing.T) {
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
				data, _ := proto.Marshal(decodePost(post))
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicName,
					Value: data,
					Key:   post.ID.String(),
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
				data, _ := proto.Marshal(decodePost(post))
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicName,
					Value: data,
					Key:   post.ID.String(),
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
			err := p.Send(tt.args.ctx, tt.args.post)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
