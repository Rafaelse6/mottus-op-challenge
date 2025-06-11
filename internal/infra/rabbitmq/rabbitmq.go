package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewRabbitMQChannel(url string) (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		"motos", // queue name
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	closeFunc := func() {
		ch.Close()
		conn.Close()
	}

	return ch, closeFunc, nil
}
