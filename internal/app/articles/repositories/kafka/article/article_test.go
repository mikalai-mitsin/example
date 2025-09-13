package repositories

import (
	"context"
	"errors"
	"reflect"
	"testing"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewArticleEventProducer(t *testing.T) {
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
		want *ArticleEventProducer
	}{
		{
			name: "ok",
			args: args{
				producer: mockProducer,
				logger:   mockLogger,
			},
			want: &ArticleEventProducer{
				producer: mockProducer,
				logger:   mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArticleEventProducer(tt.args.producer, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewArticleEventProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleEventProducer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockProducer := NewMockproducer(ctrl)
	ctx := context.Background()
	article := entities.NewMockArticle(t)
	type fields struct {
		producer producer
		logger   logger
	}
	type args struct {
		ctx     context.Context
		article entities.Article
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
				ctx:     ctx,
				article: article,
			},
			setup: func() {
				data, _ := proto.Marshal(decodeArticle(article))
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicName,
					Value: data,
					Key:   article.ID.String(),
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
				ctx:     ctx,
				article: article,
			},
			setup: func() {
				data, _ := proto.Marshal(decodeArticle(article))
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicName,
					Value: data,
					Key:   article.ID.String(),
				}).Return(errors.New("test error"))
			},
			wantErr: errs.FromKafkaError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			p := &ArticleEventProducer{
				producer: tt.fields.producer,
				logger:   tt.fields.logger,
			}
			err := p.Send(tt.args.ctx, tt.args.article)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
