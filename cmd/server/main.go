package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/controller"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/infra/rabbitmq"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/repository"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/service"
)

func main() {
	ch, closeConn, err := rabbitmq.NewRabbitMQChannel("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar com RabbitMQ: %v", err)
	}
	defer closeConn()

	err = rabbitmq.StartConsumer(ch, "motos")
	if err != nil {
		log.Fatalf("Erro ao iniciar consumidor RabbitMQ: %v", err)
	}

	repo := repository.NewInMemoryMotoRepository()
	publisher := rabbitmq.NewRabbitMQPublisher(ch)
	svc := service.NewMotoService(repo, publisher)
	ctrl := controller.NewMotoController(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/motos", ctrl.Create)

	log.Println("Servidor iniciado em http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
