package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMotoDTO(t *testing.T) {
	t.Run("should create valid MotoDTO", func(t *testing.T) {
		dto := NewMotoDTO(2024, "Honda CB", "ABC-1234")

		assert.Equal(t, 2024, dto.Year)
		assert.Equal(t, "Honda CB", dto.Model)
		assert.Equal(t, "ABC-1234", dto.Plate)
	})
}

func TestMotoDTO_ToEntity(t *testing.T) {
	tests := []struct {
		name    string
		dto     *MotoDTO
		wantErr bool
	}{
		{
			name:    "valid DTO",
			dto:     NewMotoDTO(2024, "Valid Model", "VALID-123"),
			wantErr: false,
		},
		{
			name:    "invalid year",
			dto:     NewMotoDTO(0, "Invalid Year", "INVL-123"),
			wantErr: true,
		},
		{
			name:    "empty plate",
			dto:     NewMotoDTO(2024, "Empty Plate", ""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dto.ToEntity()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMotoDTO_ToEntity_Validation(t *testing.T) {
	t.Run("should convert to entity with same values", func(t *testing.T) {
		dto := NewMotoDTO(2024, "Test Model", "TEST-123")
		moto, err := dto.ToEntity()

		assert.NoError(t, err)
		assert.Equal(t, dto.Year, moto.Year)
		assert.Equal(t, dto.Model, moto.Model)
		assert.Equal(t, dto.Plate, moto.Plate)
	})
}
