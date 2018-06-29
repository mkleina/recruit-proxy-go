package rabbitclient

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type RabbitClient struct {
	user string
	pass string
	host string
	port int
}

func makeQueue(channel *amqp.Channel, name string, exclusive bool) (*amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		name,
		false,
		false,
		exclusive,
		false,
		nil,
	)
	return &queue, err
}

func consume(channel *amqp.Channel, queueName string, exclusive bool) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		queueName,
		"recruit-proxy-go",
		true,
		exclusive,
		false,
		true,
		nil,
	)
}

func publish(channel *amqp.Channel, exchange string, routingKey string, replyTo string, correlationID string, body string) error {
	return channel.Publish(
		exchange,   // Exchange
		routingKey, // Routing key
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          []byte(body),
			ReplyTo:       replyTo,
			CorrelationId: correlationID,
		},
	)
}

// NewClient created new RabbitMQ / AMQP client
func NewClient(user string, pass string, host string, port int) *RabbitClient {
	return &RabbitClient{
		user: user,
		pass: pass,
		host: host,
		port: port,
	}
}

// GetReply sends request to AMQP queue with provided messageBody and returns response
func (r *RabbitClient) GetReply(queueName string, messageBody string) ([]byte, error) {
	correlationID := uuid.Must(uuid.NewRandom()).String()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", r.user, r.pass, r.host, r.port))
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		return nil, err
	}

	// Ensure that publish queue exists
	publishQueue, err := makeQueue(ch, queueName, false)
	if err != nil {
		return nil, err
	}

	// Ensure that reply queue exists and can be used exclusively
	replyQueue, err := makeQueue(ch, "", true)
	if err != nil {
		return nil, err
	}

	// Start consuming messages from reply queue
	msgs, err := consume(ch, replyQueue.Name, true)
	if err != nil {
		return nil, err
	}

	err = publish(ch, "", publishQueue.Name, replyQueue.Name, correlationID, messageBody)
	if err != nil {
		return nil, err
	}

	for msg := range msgs {
		if msg.CorrelationId == correlationID {
			return msg.Body, nil
		}
	}
	return nil, nil
}
