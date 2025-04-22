package service

import (
	"errors"
	"testing"
	"time"

	"github.com/snavarro/microtracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPackageRepository is a mock implementation of domain.PackageRepository
type MockPackageRepository struct {
	mock.Mock
}

func (m *MockPackageRepository) FindByID(id string) (*domain.Package, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Package), args.Error(1)
}

func (m *MockPackageRepository) FindAll(page, size int) ([]domain.Package, int64, error) {
	args := m.Called(page, size)
	return args.Get(0).([]domain.Package), args.Get(1).(int64), args.Error(2)
}

func (m *MockPackageRepository) Search(query string, page, size int) ([]domain.Package, int64, error) {
	args := m.Called(query, page, size)
	return args.Get(0).([]domain.Package), args.Get(1).(int64), args.Error(2)
}

func (m *MockPackageRepository) Create(pkg *domain.Package) error {
	args := m.Called(pkg)
	return args.Error(0)
}

func (m *MockPackageRepository) Update(pkg *domain.Package) error {
	args := m.Called(pkg)
	return args.Error(0)
}

func (m *MockPackageRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestPackageService_GetPackage(t *testing.T) {
	mockRepo := new(MockPackageRepository)
	service := NewPackageService(mockRepo)

	t.Run("successful get", func(t *testing.T) {
		expectedPkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("FindByID", "123").Return(expectedPkg, nil)

		pkg, err := service.GetPackage("123")

		assert.NoError(t, err)
		assert.Equal(t, expectedPkg, pkg)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("FindByID", "456").Return(nil, errors.New("package not found"))

		pkg, err := service.GetPackage("456")

		assert.Error(t, err)
		assert.Nil(t, pkg)
		mockRepo.AssertExpectations(t)
	})
}

func TestPackageService_ListPackages(t *testing.T) {
	mockRepo := new(MockPackageRepository)
	service := NewPackageService(mockRepo)

	t.Run("successful list", func(t *testing.T) {
		expectedPackages := []domain.Package{
			{
				PackageID: "123",
				Sender: domain.Address{
					Name:    "John Doe",
					Address: "123 Main St",
				},
			},
			{
				PackageID: "456",
				Sender: domain.Address{
					Name:    "Jane Doe",
					Address: "456 Oak St",
				},
			},
		}

		mockRepo.On("FindAll", 1, 10).Return(expectedPackages, int64(2), nil)

		packages, total, err := service.ListPackages(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedPackages, packages)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo.On("FindAll", 1, 10).Return([]domain.Package{}, int64(0), errors.New("database error"))

		packages, total, err := service.ListPackages(1, 10)

		assert.Error(t, err)
		assert.Empty(t, packages)
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})
}

func TestPackageService_CreatePackage(t *testing.T) {
	mockRepo := new(MockPackageRepository)
	service := NewPackageService(mockRepo)

	t.Run("successful create", func(t *testing.T) {
		pkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
			Recipient: domain.Address{
				Name:    "Jane Doe",
				Address: "456 Oak St",
			},
			Origin:        "New York",
			Destination:   "Los Angeles",
			CurrentStatus: "created",
		}

		mockRepo.On("Create", pkg).Return(nil)

		err := service.CreatePackage(pkg)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		pkg := &domain.Package{
			PackageID: "456",
		}

		err := service.CreatePackage(pkg)

		assert.Error(t, err)
		assert.Equal(t, "recipient name and address are required", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestPackageService_UpdatePackage(t *testing.T) {
	mockRepo := new(MockPackageRepository)
	service := NewPackageService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		pkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
			Recipient: domain.Address{
				Name:    "Jane Doe",
				Address: "456 Oak St",
			},
			Origin:        "New York",
			Destination:   "Los Angeles",
			CurrentStatus: "updated",
		}

		mockRepo.On("Update", pkg).Return(nil)

		err := service.UpdatePackage(pkg)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		pkg := &domain.Package{
			PackageID: "456",
		}

		err := service.UpdatePackage(pkg)

		assert.Error(t, err)
		assert.Equal(t, "recipient name and address are required", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestPackageService_DeletePackage(t *testing.T) {
	mockRepo := new(MockPackageRepository)
	service := NewPackageService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		mockRepo.On("Delete", "123").Return(nil)

		err := service.DeletePackage("123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo.On("Delete", "456").Return(errors.New("database error"))

		err := service.DeletePackage("456")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
