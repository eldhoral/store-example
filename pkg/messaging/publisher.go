package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"strings"
)

// Always send to partition 0
type PartitionZero struct {
}

func (cp PartitionZero) Balance(msg kafka.Message, partitions ...int) (partition int) {
	return 0
}

// NewKafkaPublisher create kafka publisher. Ref: https://github.com/segmentio/kafka-go#writing-to-multiple-topics
// Mixpanel and MoEngage as well
func NewKafkaPublisher(brokers string, topic string) *kafkaPublisher {
	w := &kafka.Writer{
		Addr: kafka.TCP(strings.Split(brokers, ",")...),
		// NOTE: When Topic is not defined here, each Message must define it instead.
		Balancer: &PartitionZero{},
	}
	seg := &kafkaPublisher{writer: w, currentTopic: topic, brokers: brokers}

	// Check topic exists before init
	if brokers != "" {
		logrus.Infoln("Load Kafka service - ", topic)
		seg.CreateTopicIfNotExist(topic)
	} else {
		logrus.Warning("Brokers not set. Check ULMS_KAFKA_BROKERS .env")
	}

	return seg
}

type KafkaPublisher interface {
	Publish(ctx context.Context, data interface{}, topic string) error
	CreateTopicIfNotExist(topic string) bool
}

type kafkaPublisher struct {
	writer       *kafka.Writer
	currentTopic string
	brokers      string
}

func (p *kafkaPublisher) Publish(ctx context.Context, data interface{}, topic string) error {
	p.checkTopicExist(topic)

	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "error on preparing kafka message")
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Key:   []byte("Key"),
			Value: bytes,
		},
	)
	fmt.Println("[INFO] Publish to topic:", topic, "- Data", data, err)
	return err
}

func (p *kafkaPublisher) CreateTopicIfNotExist(topic string) bool {

	hosts := strings.Split(p.brokers, ",")
	for _, host := range hosts {
		_, err := kafka.DialLeader(context.Background(), "tcp", host, topic, 0)

		if err != nil {
			logrus.Error(err)
			return false
		}
	}

	logrus.Infoln("Check topic:", topic)
	return true
}

func (p *kafkaPublisher) checkTopicExist(topic string) {
	if p.currentTopic != topic && p.CreateTopicIfNotExist(topic) {
		p.currentTopic = topic
	}
}
