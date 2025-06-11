package repository

import (
	"sync"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/google/uuid"
)

type InMemoryMotoRepository struct {
	motos map[uuid.UUID]*entity.Moto
	mu    sync.RWMutex
}

func NewInMemoryMotoRepository() *InMemoryMotoRepository {
	return &InMemoryMotoRepository{
		motos: make(map[uuid.UUID]*entity.Moto),
	}
}

func (r *InMemoryMotoRepository) Save(moto *entity.Moto) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.motos[moto.ID] = moto
	return nil
}

func (r *InMemoryMotoRepository) FindByID(id uuid.UUID) (*entity.Moto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if moto, ok := r.motos[id]; ok {
		return moto, nil
	}
	return nil, nil
}

func (r *InMemoryMotoRepository) FindByPlate(plate string) (*entity.Moto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, moto := range r.motos {
		if moto.Plate == plate {
			return moto, nil
		}
	}
	return nil, nil
}

func (r *InMemoryMotoRepository) UpdatePlate(id uuid.UUID, newPlate string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if moto, ok := r.motos[id]; ok {
		moto.Plate = newPlate
	}
}

func (r *InMemoryMotoRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.motos, id)
	return nil
}
