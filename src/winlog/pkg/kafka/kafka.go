package kafka

import (
	"context"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/prometheus/common/log"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/prom"
)

type Kafka struct {
	client  sarama.ConsumerGroup
	stopped chan struct{}
	msgs    chan []byte
	Name    string
	Version string
	Oldest  bool
	Brokers string
	Group   string
	Topic   string

	Sasl struct {
		Username string
		Password string
	}

	consumer *Consumer
	context  context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// Sarama consumer group
type Consumer struct {
	Name  string
	ready chan bool
	msgs  chan []byte
}

type ConsumerError struct {
	UnixTime int64
	Error    error
}

// NewKafka get instance of kafka reader
func NewKafka() *Kafka {
	return &Kafka{}
}

// Init Initialise the kafka instance with configuration
func (k *Kafka) Init() error {
	k.msgs = make(chan []byte, 1024)
	k.stopped = make(chan struct{})
	k.consumer = &Consumer{
		Name:  k.Name,
		msgs:  k.msgs,
		ready: make(chan bool),
	}
	k.context, k.cancel = context.WithCancel(context.Background())
	return nil
}

// Msgs returns the message from kafka
func (k *Kafka) Msgs() chan []byte {
	return k.msgs
}

// Start kafka consumer, uses saram library
func (k *Kafka) Start() error {
	config := sarama.NewConfig()

	if k.Version != "" {
		version, err := sarama.ParseKafkaVersion(k.Version)
		if err != nil {
			return err
		}
		config.Version = version
	}
	// check for authentication
	if k.Sasl.Username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = k.Sasl.Username
		config.Net.SASL.Password = k.Sasl.Password
	}
	if k.Oldest { // set to read the oldest message from last commit point
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	log.Info("start to dial kafka ", k.Brokers)
	client, err := sarama.NewConsumerGroup(strings.Split(k.Brokers, ","), k.Group, config)
	if err != nil {
		return err
	}

	k.client = client

	k.wg.Add(1)
	go func() {
		defer k.wg.Done()

		for {
			if err := k.client.Consume(k.context, strings.Split(k.Topic, ","), k.consumer); err != nil {
				prom.KafkaConsumerErrors.WithLabelValues(k.Topic, k.Group).Inc()
				log.Error("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if k.context.Err() != nil {
				return
			}
			k.consumer.ready = make(chan bool)
		}
	}()

	<-k.consumer.ready
	return nil
}

// Stop kafka consumer and close all connections
func (k *Kafka) Stop() error {
	k.cancel()
	k.wg.Wait()

	_ = k.client.Close()
	close(k.msgs)
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {

	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		consumer.msgs <- message.Value
		session.MarkMessage(message, "")
		prom.KafkaConsumerTotal.WithLabelValues(consumer.Name, message.Topic).Inc()

	}

	return nil
}
