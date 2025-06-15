package repository

import (
	"errors"
	"fmt"
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

	if moto == nil {
		return fmt.Errorf("moto cannot be nil")
	}

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

func (r *InMemoryMotoRepository) UpdatePlate(id uuid.UUID, newPlate string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	moto, exists := r.motos[id]
	if !exists {
		return errors.New("moto not found")
	}

	if moto, ok := r.motos[id]; ok {
		moto.Plate = newPlate
	}

	for _, m := range r.motos {
		if m.Plate == newPlate && m.ID != id {
			return errors.New("plate already registered")
		}
	}

	moto.Plate = newPlate
	return nil
}

func (r *InMemoryMotoRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.motos[id]; !exists {
		return errors.New("moto not found")
	}

	delete(r.motos, id)
	return nil
}

func (r *InMemoryMotoRepository) List(plateFilter string) ([]*entity.Moto, error) {
	r.mu.Lock()
	defer r.mu.RUnlock()

	var motos []*entity.Moto
	for _, moto := range r.motos {
		if plateFilter == "" || moto.Plate == plateFilter {
			motos = append(motos, moto)
		}
	}

	return motos, nil
}
