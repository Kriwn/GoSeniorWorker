package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultQReqName = "qReq"
	defaultQResName = "qRes"
)

type MessageHandler func(context.Context, amqp.Delivery) error

type RabbitmqService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	QReq    amqp.Queue
	QRes    amqp.Queue
}

// TODO pass conn to this
func NewRabbitmqService() (*RabbitmqService, error) {
	connURL := os.Getenv("RABBITMQ_URL")
	if connURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}

	qReqName := os.Getenv("QUEUE_NAME")
	if qReqName == "" {
		qReqName = defaultQReqName
	}

	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	qReq, err := ch.QueueDeclare(
		qReqName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to declare qReq queue %q: %w", qReqName, err)
	}


	service := &RabbitmqService{
		Conn:    conn,
		Channel: ch,
		QReq:    qReq,
	}

	log.Printf("RabbitMQ connected. consume=%s publish=%s", service.QReq.Name, service.QRes.Name)
	return service, nil
}

func (r *RabbitmqService) Consume(ctx context.Context, handler MessageHandler) error {

	fmt.Println("START CONSUMER:", r.QReq.Name)

	if err := r.Channel.Qos(1, 0, false); err != nil {
		return err
	}

	err := r.Channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return err
	}


	msgs, err := r.Channel.Consume(
		r.QReq.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			if err := handler(ctx, msg); err != nil {
				fmt.Printf("Message handling failed: %v\n", err)
				// msg.Nack(false)
				break
			}

			msg.Ack(false)
		}
	}
}

func (r *RabbitmqService) CreatePublish(ctx context.Context, body []byte) error {
	if r == nil || r.Channel == nil {
		return fmt.Errorf("rabbitmq service is not initialized")
	}

	return r.Channel.PublishWithContext(
		ctx,
		"",
		r.QRes.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (r *RabbitmqService) Close() error {
	if r == nil {
		return nil
	}

	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			return err
		}
	}

	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
