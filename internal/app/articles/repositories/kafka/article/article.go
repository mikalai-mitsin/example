package repositories

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"google.golang.org/protobuf/proto"
)

const (
	topicName = "example.articles.article.v1"
)

type ArticleEventProducer struct {
	producer producer
	logger   logger
}

func NewArticleEventProducer(
	producer producer,
	logger logger,
) *ArticleEventProducer {
	return &ArticleEventProducer{producer: producer, logger: logger}
}

func (p *ArticleEventProducer) Send(ctx context.Context, article entities.Article) error {
	data, err := proto.Marshal(decodeArticle(article))
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicName,
		Value: data,
		Key:   article.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}
