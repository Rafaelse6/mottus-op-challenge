package service

import (
	"testing"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks

type MockMotoRepository struct {
	mock.Mock
}

func (m *MockMotoRepository) Save(moto *entity.Moto) error {
	args := m.Called(moto)
	return args.Error(0)
}

func (m *MockMotoRepository) FindByID(id uuid.UUID) (*entity.Moto, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Moto), args.Error(1)
}

func (m *MockMotoRepository) FindByPlate(plate string) (*entity.Moto, error) {
	args := m.Called(plate)
	if moto, ok := args.Get(0).(*entity.Moto); ok {
		return moto, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMotoRepository) UpdatePlate(id uuid.UUID, newPlate string) {
	m.Called(id, newPlate)
}

func (m *MockMotoRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(queue string, body []byte) error {
	args := m.Called(queue, body)
	return args.Error(0)
}

// Tests

func TestCreateMoto_Success(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)

	service := NewMotoService(mockRepo, mockPublisher)

	mockRepo.On("FindByPlate", "ABC-1234").Return(nil, nil)
	mockRepo.On("Save", mock.Anything).Return(nil)
	mockPublisher.On("Publish", "motos", mock.Anything).Return(nil)

	moto, err := service.CreateMoto(2024, "Honda", "ABC-1234")

	assert.NoError(t, err)
	assert.NotNil(t, moto)
	assert.Equal(t, "Honda", moto.Model)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestCreateMoto_PlateAlreadyExists(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)

	service := NewMotoService(mockRepo, mockPublisher)

	existingMoto := &entity.Moto{ID: uuid.New(), Year: 2023, Model: "Yamaha", Plate: "XYZ-9876"}
	mockRepo.On("FindByPlate", "XYZ-9876").Return(existingMoto, nil)

	moto, err := service.CreateMoto(2024, "Honda", "XYZ-9876")

	assert.ErrorIs(t, err, ErrPlateAlreadyExists)
	assert.Nil(t, moto)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything)
}

func TestCreateMoto_InvalidData(t *testing.T) {
	repo := new(MockMotoRepository)
	publisher := new(MockPublisher)
	service := NewMotoService(repo, publisher)

	repo.On("FindByPlate", mock.Anything).Return(nil, nil)

	tests := []struct {
		name    string
		year    int
		model   string
		plate   string
		wantErr bool
	}{
		{"year is zero", 0, "Model", "ABC-1234", true},
		{"plate is empty", 2024, "Model", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateMoto(tt.year, tt.model, tt.plate)
			assert.Error(t, err)
		})
	}
}
