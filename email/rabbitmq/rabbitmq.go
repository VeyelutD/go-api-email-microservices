package rabbitmq

import (
	"context"
	"fmt"
	ampq "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitClient struct {
	conn *ampq.Connection
	ch   *ampq.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*ampq.Connection, error) {
	return ampq.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *ampq.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err := ch.Confirm(false); err != nil {
		return nil, err
	}
	err = ch.ExchangeDeclare(
		"email_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) (ampq.Queue, error) {
	queue, err := rc.ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	if err != nil {
		return ampq.Queue{}, err
	}
	return queue, nil
}
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options ampq.Publishing) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx, exchange, routingKey, true, false, options)
	if err != nil {
		return err
	}
	log.Println(confirmation.Wait())
	return nil
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan ampq.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rc RabbitClient) ApplyQos(prefetchCount int, prefetchSize int, isGlobal bool) error {
	return rc.ch.Qos(prefetchCount, prefetchSize, isGlobal)
}
