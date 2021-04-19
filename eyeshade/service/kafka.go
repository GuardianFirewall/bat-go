package eyeshade

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/brave-intl/bat-go/eyeshade/avro"
	avrocontribution "github.com/brave-intl/bat-go/eyeshade/avro/contribution"
	avroreferral "github.com/brave-intl/bat-go/eyeshade/avro/referral"
	avrosettlement "github.com/brave-intl/bat-go/eyeshade/avro/settlement"
	avrosuggestion "github.com/brave-intl/bat-go/eyeshade/avro/suggestion"
	appctx "github.com/brave-intl/bat-go/utils/context"
	errorutils "github.com/brave-intl/bat-go/utils/errors"
	kafkautils "github.com/brave-intl/bat-go/utils/kafka"
	stringutils "github.com/brave-intl/bat-go/utils/string"
	"github.com/segmentio/kafka-go"
)

var (
	// KeyToEncoder maps a key to an avro topic bundle
	KeyToEncoder = map[string]avro.TopicBundle{
		avro.TopicKeys.Settlement:   avrosettlement.New(),
		avro.TopicKeys.Referral:     avroreferral.New(),
		avro.TopicKeys.Suggestion:   avrosuggestion.New(),
		avro.TopicKeys.Contribution: avrocontribution.New(),
	}
)

// BatchMessagesHandler handles many messages being consumed at once
type BatchMessagesHandler interface {
	Topic() string
}

// BatchMessageConsumer holds methods for batch method consumption
type BatchMessageConsumer interface {
	Consume(erred chan error)
	IsConsuming() bool
	IdleFor(time.Duration) bool
	Read() ([]kafka.Message, bool, error)
	Handler([]kafka.Message) error
	Commit([]kafka.Message) error
}

// BatchMessageProducer holds methods for batch method producer
type BatchMessageProducer interface {
	Produce(
		context.Context,
		avro.TopicBundle,
		...interface{},
	) error
}

// MessageHandler holds information about a single handler
type MessageHandler struct {
	key         string
	isConsuming bool
	isHandling  bool
	lastTick    time.Time
	bundle      avro.TopicBundle
	service     *Service
	reader      *kafka.Reader
	writer      *kafka.Writer
	dialer      *kafka.Dialer
}

// IsConsuming checks if the message handler is consuming
func (con *MessageHandler) IsConsuming() bool {
	return con.isConsuming
}

// IdleFor returns notes whether or not the consumer is temporarily done consuming
func (con *MessageHandler) IdleFor(duration time.Duration) bool {
	return !con.isHandling && (con.lastTick.IsZero() || time.Now().After(con.lastTick.Add(duration)))
}

// Topic returns the topic of the message handler
func (con *MessageHandler) Topic() string {
	return avro.KeyToTopic[con.bundle.Key()]
}

// Collect collects messages from the reader
func (con *MessageHandler) Collect(
	msgCh chan<- kafka.Message,
	errCh chan<- error,
) {
	msg, err := con.reader.FetchMessage(con.Context())
	if err != nil {
		errCh <- err
	} else {
		msgCh <- msg
	}
}

// Read reads messages
func (con *MessageHandler) Read() ([]kafka.Message, bool, error) {
	msg, err := con.reader.FetchMessage(con.Context())
	msgs := []kafka.Message{}
	if err != nil {
		return msgs, false, err
	}
	return append(msgs, msg), false, err
}

// Handler handles the batch of messages
func (con *MessageHandler) Handler(msgs []kafka.Message) error {
	handler := Handlers[con.bundle.Key()]
	if handler == nil {
		return errors.New("unknown handler asked for during message handle")
	}
	return handler(con, msgs)
}

// Context returns the context from the service
func (con *MessageHandler) Context() context.Context {
	return con.service.Context()
}

// Commit commits messages that have been read and inserted
func (con *MessageHandler) Commit(msgs []kafka.Message) error {
	return con.reader.CommitMessages(con.Context(), msgs...)
}

// Consume starts the consumer
func (con *MessageHandler) Consume(
	erred chan error,
) {
	var err error
	con.isConsuming = true
	for { // loop to continue consuming
		if !con.isConsuming {
			break
		}
		con.isHandling = false
		msgs, _, e := con.Read()
		if e != nil {
			err = errorutils.Wrap(e, "during read")
			break
		}
		if len(msgs) == 0 {
			continue
		}
		con.isHandling = true
		con.lastTick = time.Now()
		e = con.Handler(msgs)
		if e != nil {
			err = errorutils.Wrap(e, "during handler")
			break
		}
		e = con.Commit(msgs)
		if e != nil {
			err = errorutils.Wrap(e, "during commit")
			break
		}
		con.lastTick = time.Now()
	}
	con.isHandling = false
	con.isConsuming = false
	if err != nil {
		erred <- errorutils.Wrap(
			err,
			fmt.Sprintf("error in topic - %s", con.reader.Config().Topic),
		)
	}
}

// Produce produces messages
func (con *MessageHandler) Produce(
	ctx context.Context,
	encoder avro.TopicBundle,
	encodables ...interface{},
) error {
	if len(encodables) == 0 {
		return nil
	}
	messages, err := encoder.ManyToBinary(encodables...)
	if err != nil {
		return err
	}
	return con.writer.WriteMessages(
		ctx,
		*messages...,
	)
}

// WithTopicAutoCreation creates topics used by producers and consumers
func WithTopicAutoCreation(service *Service) error {
	keys := map[string]int{}
	if service.producers != nil {
		for key := range service.producers {
			keys[key] = 0
		}
	}
	if service.consumers != nil {
		for key := range service.consumers {
			keys[key] = 0
		}
	}

	ctx := service.Context()
	kafkaBrokers := ctx.Value(appctx.KafkaBrokersCTXKey).(string)
	broker := stringutils.SplitAndTrim(kafkaBrokers)[0]
	conn, err := service.dialer.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("connection close failed", err)
		}
	}()
	reps, err := strconv.Atoi(os.Getenv("KAFKA_REPLICATIONS"))
	if err != nil {
		return err
	}
	partitions, err := strconv.Atoi(os.Getenv("KAFKA_PARTITIONS"))
	if err != nil {
		return err
	}
	for key := range keys {
		topic := avro.KeyToTopic[key]
		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: reps,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// WithConsumer sets up a consumer on the service
func WithConsumer(
	topicKeys ...string,
) func(*Service) error {
	return func(service *Service) error {
		if service.consumers == nil {
			// can't access consumer keys until make is called
			service.consumers = make(map[string]BatchMessageConsumer)
		}
		for _, topicKey := range topicKeys {
			reader, dialer, err := kafkautils.InitKafkaReader(
				service.Context(),
				avro.KeyToTopic[topicKey],
				service.dialer,
			)
			if err != nil {
				return err
			}
			consumer := &MessageHandler{
				key:     topicKey,
				bundle:  KeyToEncoder[topicKey],
				reader:  reader,
				dialer:  dialer,
				service: service,
			}
			service.dialer = dialer
			service.consumers[topicKey] = BatchMessageConsumer(consumer)
		}
		return nil
	}
}

// Consumer retrieves a consumer off of the service
func (service *Service) Consumer(key string) BatchMessageConsumer {
	return service.consumers[key]
}

// WithProducer sets up a consumer on the service
func WithProducer(
	topicKeys ...string,
) func(*Service) error {
	return func(service *Service) error {
		if service.producers == nil {
			// can't access producer keys until make is called
			service.producers = make(map[string]BatchMessageProducer)
		}
		for _, topicKey := range topicKeys {
			writer, dialer, err := kafkautils.InitKafkaWriter(
				service.Context(),
				avro.KeyToTopic[topicKey],
				service.dialer,
			)
			if err != nil {
				return err
			}
			producer := &MessageHandler{
				key:     topicKey,
				bundle:  KeyToEncoder[topicKey],
				writer:  writer,
				dialer:  dialer,
				service: service,
			}
			service.dialer = dialer
			service.producers[topicKey] = BatchMessageProducer(producer)
		}
		return nil
	}
}