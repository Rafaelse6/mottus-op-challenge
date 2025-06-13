package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func NewRabbitMQChannel(amqpURL string) (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// Declara a fila (caso ainda não exista)
	_, err = ch.QueueDeclare(
		"motos",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	log.Println("RabbitMQ conectado com sucesso")

	closeFunc := func() {
		ch.Close()
		conn.Close()
	}

	return ch, closeFunc, nil
}
