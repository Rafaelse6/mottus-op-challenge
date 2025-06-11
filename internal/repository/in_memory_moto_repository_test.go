package repository

import (
	"testing"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryMotoRepository_Save(t *testing.T) {
	repo := NewInMemoryMotoRepository()
	moto := &entity.Moto{
		ID:    uuid.New(),
		Year:  2024,
		Model: "Test Model",
		Plate: "TEST-123",
	}

	t.Run("should save moto successfully", func(t *testing.T) {
		err := repo.Save(moto)
		assert.NoError(t, err)

		// Verifica se a moto foi salva
		repo.mu.RLock()
		defer repo.mu.RUnlock()
		savedMoto, exists := repo.motos[moto.ID]
		assert.True(t, exists)
		assert.Equal(t, moto, savedMoto)
	})

	t.Run("should return error when moto is nil", func(t *testing.T) {
		err := repo.Save(nil)
		assert.Error(t, err)
		assert.Equal(t, "moto cannot be nil", err.Error())
	})
}

func TestInMemoryMotoRepository_FindByID(t *testing.T) {
	repo := NewInMemoryMotoRepository()
	moto := &entity.Moto{
		ID:    uuid.New(),
		Year:  2024,
		Model: "Test Model",
		Plate: "TEST-123",
	}

	repo.Save(moto)

	t.Run("should find moto by ID", func(t *testing.T) {
		foundMoto, err := repo.FindByID(moto.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundMoto)
		assert.Equal(t, moto.ID, foundMoto.ID)
	})

	t.Run("should return nil when moto not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		foundMoto, err := repo.FindByID(nonExistentID)
		assert.NoError(t, err)
		assert.Nil(t, foundMoto)
	})
}

func TestInMemoryMotoRepository_FindByPlate(t *testing.T) {
	repo := NewInMemoryMotoRepository()
	moto := &entity.Moto{
		ID:    uuid.New(),
		Year:  2024,
		Model: "Test Model",
		Plate: "TEST-123",
	}

	repo.Save(moto)

	t.Run("should find moto by plate", func(t *testing.T) {
		foundMoto, err := repo.FindByPlate("TEST-123")
		assert.NoError(t, err)
		assert.NotNil(t, foundMoto)
		assert.Equal(t, moto.ID, foundMoto.ID)
	})

	t.Run("should return nil when plate not found", func(t *testing.T) {
		foundMoto, err := repo.FindByPlate("NON-EXISTENT")
		assert.NoError(t, err)
		assert.Nil(t, foundMoto)
	})
}

func TestInMemoryMotoRepository_UpdatePlate(t *testing.T) {
	repo := NewInMemoryMotoRepository()
	moto := &entity.Moto{
		ID:    uuid.New(),
		Year:  2024,
		Model: "Test Model",
		Plate: "TEST-123",
	}

	repo.Save(moto)

	t.Run("should update plate successfully", func(t *testing.T) {
		newPlate := "UPDATED-456"
		repo.UpdatePlate(moto.ID, newPlate)

		foundMoto, err := repo.FindByID(moto.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundMoto)
		assert.Equal(t, newPlate, foundMoto.Plate)
	})
}

func TestInMemoryMotoRepository_Delete(t *testing.T) {
	repo := NewInMemoryMotoRepository()
	moto := &entity.Moto{
		ID:    uuid.New(),
		Year:  2024,
		Model: "Test Model",
		Plate: "TEST-123",
	}

	repo.Save(moto)

	t.Run("should delete moto successfully", func(t *testing.T) {
		err := repo.Delete(moto.ID)
		assert.NoError(t, err)

		foundMoto, err := repo.FindByID(moto.ID)
		assert.NoError(t, err)
		assert.Nil(t, foundMoto)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.Delete(nonExistentID)
		assert.Error(t, err)
		assert.Equal(t, "moto not found", err.Error())
	})
}
