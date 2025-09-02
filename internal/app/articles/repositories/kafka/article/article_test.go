package events

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/IBM/sarama"
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func TestArticleEventProducer_Created(t *testing.T) {
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
				data, _ := json.Marshal(article)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventCreated,
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
				data, _ := json.Marshal(article)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventCreated,
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
			err := p.Created(tt.args.ctx, tt.args.article)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestArticleEventProducer_Updated(t *testing.T) {
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
				data, _ := json.Marshal(article)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventUpdated,
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
				data, _ := json.Marshal(article)
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventUpdated,
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
			err := p.Updated(tt.args.ctx, tt.args.article)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestArticleEventProducer_Deleted(t *testing.T) {
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
				id:  article.ID,
			},
			setup: func() {
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventDeleted,
					Value: sarama.ByteEncoder(article.ID.String()),
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
				ctx: ctx,
				id:  article.ID,
			},
			setup: func() {
				mockProducer.EXPECT().Send(gomock.Any(), &kafka.Message{
					Topic: topicEventDeleted,
					Value: []byte(article.ID.String()),
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
			err := p.Deleted(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
