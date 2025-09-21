package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type Message struct {
	Topic string
	Value []byte
	Key   string
}
type Producer struct {
	config   *Config
	producer sarama.SyncProducer
	logger   log.Logger
}

func NewProducer(cfg *Config, logger log.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Version = sarama.V2_1_0_0
	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, errs.NewUnexpectedBehaviorError("cant build kafka producer").WithCause(err)
	}
	return &Producer{config: cfg, producer: producer, logger: logger}, nil
}
func (p *Producer) Send(_ context.Context, message *Message) error {
	msg := &sarama.ProducerMessage{
		Topic: message.Topic,
		Key:   sarama.StringEncoder(message.Key),
		Value: sarama.ByteEncoder(message.Value),
	}
	_, _, err := p.producer.SendMessage(msg)
	return err
}
func (p *Producer) Start(ctx context.Context) error {
	return nil
}
func (p *Producer) Stop(ctx context.Context) error {
	return p.producer.Close()
}
