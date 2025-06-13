package rabbitmq

import (
	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	ch *amqp.Channel
}

func NewRabbitMQPublisher(ch *amqp.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{ch: ch}
}

func (p *RabbitMQPublisher) Publish(queue string, body []byte) error {
	return p.ch.Publish(
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
