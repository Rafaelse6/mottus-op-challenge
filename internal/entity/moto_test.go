package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMoto(t *testing.T) {
	moto, err := NewMoto(2025, "Kawazaki Ninja", "123jkl")
	assert.Nil(t, err)
	assert.NotNil(t, moto)
	assert.NotEmpty(t, moto.ID)
	assert.NotEmpty(t, moto.Year)
	assert.NotEmpty(t, moto.Model)
	assert.NotEmpty(t, moto.Plate)
	assert.Equal(t, "Kawazaki Ninja", moto.Model)
	assert.Equal(t, "123jkl", moto.Plate)
}

func TestNewMoto_InvalidYear(t *testing.T) {
	_, err := NewMoto(0, "Model X", "ABC-1234")
	assert.Error(t, err)
	assert.Equal(t, "year must be positive", err.Error())
}

func TestNewMoto_EmptyPlate(t *testing.T) {
	_, err := NewMoto(2024, "Model X", "")
	assert.Error(t, err)
	assert.Equal(t, "plate cannot be empty", err.Error())
}
