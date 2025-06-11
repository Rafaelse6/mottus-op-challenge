package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/event"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/repository"
	"github.com/streadway/amqp"
)

var (
	ErrIDIsRequired       = errors.New("id is required")
	ErrInvaldiId          = errors.New("invalid id")
	ErrYearIsRequired     = errors.New("year is required")
	ErrModelIsRequired    = errors.New("model is required")
	ErrInvalidPlate       = errors.New("invalid plate")
	ErrPlateAlreadyExists = errors.New("plate already registered")
	ErrMotoNotFound       = errors.New("moto not found")
)

type MotoService struct {
	repo      repository.MotoRepository
	publisher event.Publisher
	rmqChan   *amqp.Channel
}

func NewMotoService(repo repository.MotoRepository, rmqChan *amqp.Channel) *MotoService {
	return &MotoService{
		repo:    repo,
		rmqChan: rmqChan,
	}
}

func (s *MotoService) CreateMoto(year int, model, plate string) (*entity.Moto, error) {
	existing, _ := s.repo.FindByPlate(plate)
	if existing != nil {
		return nil, ErrPlateAlreadyExists
	}

	moto, err := entity.NewMoto(year, model, plate)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(moto)
	if err != nil {
		return nil, err
	}

	// Publica o evento no RabbitMQ
	payload, err := json.Marshal(map[string]interface{}{
		"event":   "Moto created",
		"moto_id": moto.ID,
		"year":    moto.Year,
		"model":   moto.Model,
		"plate":   moto.Plate,
	})
	if err != nil {
		log.Printf("Erro ao serializar mensagem: %v", err)
		return moto, nil // Não falha o cadastro por erro no evento
	}

	err = s.rmqChan.Publish(
		"",      // exchange
		"motos", // queue name (a mesma declarada no NewRabbitMQChannel)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		log.Printf("Erro ao publicar mensagem: %v", err)
	}

	log.Printf("Moto criada e evento publicado: %s", payload)
	return moto, nil
}
