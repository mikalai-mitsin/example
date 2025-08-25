package kafka

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type Handler struct {
	Topic   string
	GroupID string
	Handler sarama.ConsumerGroupHandler
	group   sarama.ConsumerGroup
}
type Consumer struct {
	config   *Config
	client   sarama.Client
	handlers map[string]Handler
	logger   *log.Log
}

func NewConsumer(cfg *Config, logger *log.Log) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	client, err := sarama.NewClient(cfg.Brokers, config)
	if err != nil {
		return nil, errs.NewUnexpectedBehaviorError("cant build kafka client").WithCause(err)
	}
	return &Consumer{
		config:   cfg,
		handlers: make(map[string]Handler),
		client:   client,
		logger:   logger,
	}, nil
}
func (c *Consumer) AddHandler(handler Handler) {
	c.handlers[handler.GroupID] = handler
}
func (c *Consumer) Start(ctx context.Context) error {
	logger := c.logger
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	for id, handler := range c.handlers {
		consumerGroup, err := sarama.NewConsumerGroupFromClient(handler.GroupID, c.client)
		if err != nil {
			return errs.NewUnexpectedBehaviorError("cant build kafka consumer").WithCause(err)
		}
		defer consumerGroup.Close()
		handler.group = consumerGroup
		c.handlers[id] = handler
	}
	for {
		select {
		case <-sigterm:
			return nil
		default:
			for _, handler := range c.handlers {
				if err := handler.group.Consume(context.Background(), []string{handler.Topic}, handler.Handler); err != nil {
					logger.Error(
						"consume error",
						log.Error(err),
						log.String("group", handler.GroupID),
						log.String("topic", handler.Topic),
					)
				}
			}
		}
	}
}
func (c *Consumer) Stop(ctx context.Context) error {
	return c.client.Close()
}
