package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/infra/rabbitmq"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/infra/repository"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/service"
)

func main() {
	// Conectar ao banco de dados
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=1234567 dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Conectar ao RabbitMQ
	conn, err := amqp.Dial("amqp://YOUR_USERNAME:YOUR_PASSWORD@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Injetar dependências
	motoRepo := repository.NewMotoPostgresRepository(db)
	publisher := rabbitmq.NewRabbitMQPublisher(ch)
	motoService := service.NewMotoService(motoRepo, publisher)

	// Criar router com Chi
	r := chi.NewRouter()

	// Handlers
	motoHandler := handler.NewMotoHandler(motoService)

	r.Post("/motos", motoHandler.CreateMoto)
	r.Get("/motos", motoHandler.ListMotos)
	r.Put("/motos/{id}", motoHandler.UpdateMoto)
	r.Delete("/motos/{id}", motoHandler.DeleteMoto)

	// Iniciar servidor
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
