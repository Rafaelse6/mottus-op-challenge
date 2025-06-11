// internal/infra/rabbitmq/rabbitmq.go
package rabbitmq

import (
	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	Channel *amqp.Channel
}

func NewRabbitMQPublisher(channel *amqp.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{Channel: channel}
}

func (r *RabbitMQPublisher) Publish(queue string, body []byte) error {
	_, err := r.Channel.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	return r.Channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
