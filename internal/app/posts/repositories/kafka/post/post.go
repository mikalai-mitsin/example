package repositories

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"google.golang.org/protobuf/proto"
)

const (
	topicName = "example.posts.post.v1"
)

type PostEventProducer struct {
	producer producer
	logger   logger
}

func NewPostEventProducer(
	producer producer,
	logger logger,
) *PostEventProducer {
	return &PostEventProducer{producer: producer, logger: logger}
}

func (p *PostEventProducer) Send(ctx context.Context, post entities.Post) error {
	data, err := proto.Marshal(decodePost(post))
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicName,
		Value: data,
		Key:   post.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}
