package producer

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

type producer struct {
	producer sarama.SyncProducer
	enabled  bool
}

func New(brokers string, enabled bool) (*producer, error) {
	p := &producer{
		enabled: enabled,
	}

	if !enabled {
		log.Infof("kafka producer disabled for brokers %s", brokers)

		return p, nil
	}

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true

	kafkaProducer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("creating Kafka broker: %w", err)
	}

	p.producer = kafkaProducer

	return p, nil
}

func (p *producer) Produce(topic string, msg string, key string) error {
	if !p.enabled {
		return nil
	}

	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(msg),
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *producer) Close() error {
	if p.enabled {
		return p.producer.Close()
	}

	return nil
}
