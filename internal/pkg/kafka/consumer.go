package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"golang.org/x/sync/errgroup"
)

type HandlerFunc func(ctx context.Context, msg *sarama.ConsumerMessage) error
type Handler struct {
	Topic        string
	GroupID      string
	HandlerFunc  HandlerFunc
	groupHandler sarama.ConsumerGroupHandler
	group        sarama.ConsumerGroup
}

func NewHandler(topic string, groupID string, handlerFunc HandlerFunc) Handler {
	return Handler{
		Topic:        topic,
		GroupID:      groupID,
		HandlerFunc:  handlerFunc,
		groupHandler: nil,
		group:        nil,
	}
}

type Consumer struct {
	config   *Config
	client   sarama.Client
	handlers map[string]Handler
	logger   log.Logger
}

func NewConsumer(cfg *Config, logger log.Logger) (*Consumer, error) {
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
	for id, handler := range c.handlers {
		consumerGroup, err := sarama.NewConsumerGroupFromClient(handler.GroupID, c.client)
		if err != nil {
			return errs.NewUnexpectedBehaviorError("cant build kafka consumer").WithCause(err)
		}
		handler.group = consumerGroup
		handler.groupHandler = NewGroupHandler(handler.HandlerFunc, logger)
		c.handlers[id] = handler
	}
	errorgroup, ctx := errgroup.WithContext(ctx)
	for _, handler := range c.handlers {
		errorgroup.Go(func() error {
			if err := handler.group.Consume(context.Background(), []string{handler.Topic}, handler.groupHandler); err != nil {
				logger.Error(
					"consume error",
					log.Error(err),
					log.String("group", handler.GroupID),
					log.String("topic", handler.Topic),
				)
			}
			return nil
		})
	}
	return errorgroup.Wait()
}
func (c *Consumer) Stop(ctx context.Context) error {
	for _, handler := range c.handlers {
		if err := handler.group.Close(); err != nil {
			return err
		}
	}
	return c.client.Close()
}

type GroupHandler struct {
	handlerFunc HandlerFunc
	logger      log.Logger
}

func NewGroupHandler(handlerFunc HandlerFunc, logger log.Logger) *GroupHandler {
	return &GroupHandler{handlerFunc: handlerFunc, logger: logger}
}
func (h *GroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (h *GroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *GroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		ctx := context.Background()
		logger := h.logger.WithContext(ctx)
		logger.Info(
			"received message",
			log.String("topic", msg.Topic),
			log.Int32("partition", msg.Partition),
			log.Int64("offset", msg.Offset),
			log.String("key", string(msg.Key)),
			log.String("value", string(msg.Value)),
		)
		if err := h.handlerFunc(context.Background(), msg); err != nil {
			logger.Error("handled message error", log.Error(err))
			continue
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
