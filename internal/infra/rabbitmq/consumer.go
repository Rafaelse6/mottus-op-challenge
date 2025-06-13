package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func StartConsumer(ch *amqp.Channel, queue string) error {
	msgs, err := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Mensagem recebida: %s", msg.Body)
		}
	}()

	return nil
}
