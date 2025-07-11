package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/event"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/repository"
	"github.com/google/uuid"
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
}

func NewMotoService(repo repository.MotoRepository, publisher event.Publisher) *MotoService {
	return &MotoService{
		repo:      repo,
		publisher: publisher,
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

	payload, err := json.Marshal(map[string]interface{}{
		"event":   "Moto created",
		"moto_id": moto.ID,
		"year":    moto.Year,
		"model":   moto.Model,
		"plate":   moto.Plate,
	})
	if err != nil {
		log.Printf("Erro ao serializar mensagem: %v", err)
		return moto, nil
	}

	if err := s.publisher.Publish("motos", payload); err != nil {
		log.Printf("Erro ao publicar mensagem: %v", err)
	}

	log.Printf("Moto criada e evento publicado: %s", payload)
	return moto, nil
}

func (s *MotoService) ListMotos(plateFilter string) ([]*entity.Moto, error) {
	return s.repo.List(plateFilter)
}

func (s *MotoService) FindByPlate(plate string) (*entity.Moto, error) {
	if plate == "" {
		return nil, ErrInvalidPlate
	}

	moto, err := s.repo.FindByPlate(plate)
	if err != nil {
		return nil, err
	}
	if moto == nil {
		return nil, ErrMotoNotFound
	}
	return moto, nil
}

func (s *MotoService) UpdatePlate(id uuid.UUID, newPlate string) error {

	existing, _ := s.repo.FindByPlate(newPlate)
	if existing != nil && existing.ID != id {
		return ErrPlateAlreadyExists
	}

	return s.repo.UpdatePlate(id, newPlate)
}

func (s *MotoService) DeleteMoto(id uuid.UUID) error {
	return s.repo.Delete(id)
}
