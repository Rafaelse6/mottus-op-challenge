package repository

import (
	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/google/uuid"
)

type MotoRepository interface {
	Save(moto *entity.Moto) error
	FindByID(id uuid.UUID) (*entity.Moto, error)
	FindByPlate(plate string) (*entity.Moto, error)
	UpdatePlate(id uuid.UUID, newPlate string) error
	Delete(id uuid.UUID) error
	List(plateFilter string) ([]*entity.Moto, error)
}
