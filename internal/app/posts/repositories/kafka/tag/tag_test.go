package repositories

import (
	"context"
	"errors"
	"reflect"
	"testing"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewTagEventProducer(t *testing.T) {
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
		want *TagEventProducer
	}{
		{
			name: "ok",
			args: args{
				producer: mockProducer,
				logger:   mockLogger,
			},
			want: &TagEventProducer{
				producer: mockProducer,
				logger:   mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTagEventProducer(tt.args.producer, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewTagEventProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagEventProducer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	tag := entities.NewMockTag(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx context.Context
		tag entities.Tag
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
				tag: tag,
			},
			setup: func() {
				data, _ := proto.Marshal(decodeTag(tag))
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicName,
					Value: data,
					Key:   tag.ID.String(),
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
				tag: tag,
			},
			setup: func() {
				data, _ := proto.Marshal(decodeTag(tag))
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicName,
					Value: data,
					Key:   tag.ID.String(),
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &TagEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Send(tt.args.ctx, tt.args.tag)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
