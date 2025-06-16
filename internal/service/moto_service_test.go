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

func (m *MockMotoRepository) List(plateFilter string) ([]*entity.Moto, error) {
	args := m.Called(plateFilter)
	return args.Get(0).([]*entity.Moto), args.Error(1)
}

func (m *MockMotoRepository) Save(moto *entity.Moto) error {
	args := m.Called(moto)
	return args.Error(0)
}

func (m *MockMotoRepository) FindByID(id uuid.UUID) (*entity.Moto, error) {
	args := m.Called(id)

	moto, _ := args.Get(0).(*entity.Moto)
	return moto, args.Error(1)
}

func (m *MockMotoRepository) FindByPlate(plate string) (*entity.Moto, error) {
	args := m.Called(plate)
	if moto, ok := args.Get(0).(*entity.Moto); ok {
		return moto, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMotoRepository) UpdatePlate(id uuid.UUID, newPlate string) error {
	args := m.Called(id, newPlate)
	return args.Error(0)
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

func TestUpdatePlate_Success(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)
	service := NewMotoService(mockRepo, mockPublisher)

	id := uuid.New()
	newPlate := "DEF-5678"

	mockRepo.On("FindByPlate", newPlate).Return(nil, nil)
	mockRepo.On("UpdatePlate", id, newPlate).Return(nil)

	err := service.UpdatePlate(id, newPlate)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdatePlate_PlateAlreadyExists(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)
	service := NewMotoService(mockRepo, mockPublisher)

	id := uuid.New()
	existing := &entity.Moto{ID: uuid.New(), Plate: "DEF-5678"}

	mockRepo.On("FindByPlate", "DEF-5678").Return(existing, nil)

	err := service.UpdatePlate(id, "DEF-5678")
	assert.ErrorIs(t, err, ErrPlateAlreadyExists)

	mockRepo.AssertExpectations(t)
}

func TestFindByPlate_Success(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)
	service := NewMotoService(mockRepo, mockPublisher)

	expected := &entity.Moto{ID: uuid.New(), Plate: "XYZ-9999"}
	mockRepo.On("FindByPlate", "XYZ-9999").Return(expected, nil)

	result, err := service.FindByPlate("XYZ-9999")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestFindByPlate_NotFound(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)
	service := NewMotoService(mockRepo, mockPublisher)

	mockRepo.On("FindByPlate", "ABC-0000").Return(nil, nil)

	result, err := service.FindByPlate("ABC-0000")

	assert.ErrorIs(t, err, ErrMotoNotFound)
	assert.Nil(t, result)
}

func TestDeletMoto_Success(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)
	service := NewMotoService(mockRepo, mockPublisher)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(nil)

	err := service.DeleteMoto(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListMotos_WithAndWithoutFilter(t *testing.T) {
	mockRepo := new(MockMotoRepository)
	mockPublisher := new(MockPublisher)
	service := NewMotoService(mockRepo, mockPublisher)

	motos := []*entity.Moto{
		{ID: uuid.New(), Year: 2024, Model: "Honda", Plate: "ABC-1234"},
		{ID: uuid.New(), Year: 2023, Model: "Yamaha", Plate: "XYZ-5678"},
	}

	mockRepo.On("List", "").Return(motos, nil)

	result, err := service.ListMotos("")
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	filtered := []*entity.Moto{motos[1]}
	mockRepo.On("List", "XYZ-5678").Return(filtered, nil)

	resultFiltered, err := service.ListMotos("XYZ-5678")
	assert.NoError(t, err)
	assert.Len(t, resultFiltered, 1)
	assert.Equal(t, "XYZ-5678", resultFiltered[0].Plate)

	mockRepo.AssertExpectations(t)
}
