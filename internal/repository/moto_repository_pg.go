package repository

import (
	"database/sql"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/google/uuid"
)

type motoRepositoryPg struct {
	db *sql.DB
}

// Lis implements MotoRepository.
func (r *motoRepositoryPg) List(plateFilter string) ([]*entity.Moto, error) {
	panic("unimplemented")
}

// FindByID implements MotoRepository.
func (r *motoRepositoryPg) FindByID(id uuid.UUID) (*entity.Moto, error) {
	panic("unimplemented")
}

// FindByPlate implements MotoRepository.
func (r *motoRepositoryPg) FindByPlate(plate string) (*entity.Moto, error) {
	panic("unimplemented")
}

func NewMotoRepositoryPg(db *sql.DB) MotoRepository {
	return &motoRepositoryPg{db: db}
}

// UpdatePlate implements MotoRepository.
func (r *motoRepositoryPg) UpdatePlate(id uuid.UUID, newPlate string) error {
	panic("unimplemented")
}

// Delete implements MotoRepository.
func (r *motoRepositoryPg) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

func (r *motoRepositoryPg) Save(moto *entity.Moto) error {
	_, err := r.db.Exec(`INSERT INTO motos (id, year, model, plate) VALUES ($1, $2, $3, $4)`,
		moto.ID, moto.Year, moto.Model, moto.Plate)
	return err
}
