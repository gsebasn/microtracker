package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snavarro/microtracker/internal/domain"
	"github.com/snavarro/microtracker/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPackageService is a mock implementation of PackageService
type MockPackageService struct {
	mock.Mock
}

func (m *MockPackageService) GetPackage(id string) (*domain.Package, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Package), args.Error(1)
}

func (m *MockPackageService) ListPackages(page, size int) ([]domain.Package, int64, error) {
	args := m.Called(page, size)
	return args.Get(0).([]domain.Package), args.Get(1).(int64), args.Error(2)
}

func (m *MockPackageService) SearchPackages(query string, page, size int) ([]domain.Package, int64, error) {
	args := m.Called(query, page, size)
	return args.Get(0).([]domain.Package), args.Get(1).(int64), args.Error(2)
}

func (m *MockPackageService) CreatePackage(pkg *domain.Package) error {
	args := m.Called(pkg)
	return args.Error(0)
}

func (m *MockPackageService) UpdatePackage(pkg *domain.Package) error {
	args := m.Called(pkg)
	return args.Error(0)
}

func (m *MockPackageService) DeletePackage(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter(handler *PackageHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	{
		packages := api.Group("/packages")
		{
			packages.GET("", handler.ListPackages)
			packages.GET("/search", handler.SearchPackages)
			packages.GET("/:id", handler.GetPackage)
			packages.POST("", handler.CreatePackage)
			packages.PUT("/:id", handler.UpdatePackage)
			packages.DELETE("/:id", handler.DeletePackage)
		}
	}
	return router
}

func TestPackageHandler_GetPackage(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)
	router := setupTestRouter(handler)

	t.Run("success", func(t *testing.T) {
		expectedPkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetPackage", "123").Return(expectedPkg, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/packages/123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Convert response data to Package
		jsonData, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		var actualPkg domain.Package
		err = json.Unmarshal(jsonData, &actualPkg)
		assert.NoError(t, err)
		assert.Equal(t, expectedPkg.PackageID, actualPkg.PackageID)
		assert.Equal(t, expectedPkg.Sender, actualPkg.Sender)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetPackage", "456").Return(nil, service.ErrEmptyPackageID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/packages/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, service.ErrEmptyPackageID.Error(), response.Error)
	})
}

func TestPackageHandler_ListPackages(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)
	router := setupTestRouter(handler)

	t.Run("success", func(t *testing.T) {
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

		mockService.On("ListPackages", 1, 10).Return(expectedPackages, int64(2), nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/packages?page=1&size=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.Size)

		// Convert response data to []Package
		jsonData, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		var actualPackages []domain.Package
		err = json.Unmarshal(jsonData, &actualPackages)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedPackages), len(actualPackages))
		for i, expected := range expectedPackages {
			assert.Equal(t, expected.PackageID, actualPackages[i].PackageID)
			assert.Equal(t, expected.Sender, actualPackages[i].Sender)
		}
	})
}

func TestPackageHandler_SearchPackages(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)
	router := setupTestRouter(handler)

	t.Run("success", func(t *testing.T) {
		expectedPackages := []domain.Package{
			{
				PackageID: "123",
				Sender: domain.Address{
					Name:    "John Doe",
					Address: "123 Main St",
				},
			},
		}

		mockService.On("SearchPackages", "John", 1, 10).Return(expectedPackages, int64(1), nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/packages/search?query=John&page=1&size=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.Size)

		// Convert response data to []Package
		jsonData, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		var actualPackages []domain.Package
		err = json.Unmarshal(jsonData, &actualPackages)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedPackages), len(actualPackages))
		for i, expected := range expectedPackages {
			assert.Equal(t, expected.PackageID, actualPackages[i].PackageID)
			assert.Equal(t, expected.Sender, actualPackages[i].Sender)
		}
	})
}

func TestPackageHandler_CreatePackage(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)
	router := setupTestRouter(handler)

	t.Run("success", func(t *testing.T) {
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

		mockService.On("CreatePackage", pkg).Return(nil)

		jsonData, _ := json.Marshal(pkg)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Convert response data to Package
		jsonData, err = json.Marshal(response.Data)
		assert.NoError(t, err)
		var actualPkg domain.Package
		err = json.Unmarshal(jsonData, &actualPkg)
		assert.NoError(t, err)
		assert.Equal(t, pkg.PackageID, actualPkg.PackageID)
		assert.Equal(t, pkg.Sender, actualPkg.Sender)
		assert.Equal(t, pkg.Recipient, actualPkg.Recipient)
		assert.Equal(t, pkg.Origin, actualPkg.Origin)
		assert.Equal(t, pkg.Destination, actualPkg.Destination)
		assert.Equal(t, pkg.CurrentStatus, actualPkg.CurrentStatus)
	})

	t.Run("invalid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Error)
	})
}

func TestPackageHandler_UpdatePackage(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)
	router := setupTestRouter(handler)

	t.Run("success", func(t *testing.T) {
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

		mockService.On("UpdatePackage", pkg).Return(nil)

		jsonData, _ := json.Marshal(pkg)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/packages/123", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Convert response data to Package
		jsonData, err = json.Marshal(response.Data)
		assert.NoError(t, err)
		var actualPkg domain.Package
		err = json.Unmarshal(jsonData, &actualPkg)
		assert.NoError(t, err)
		assert.Equal(t, pkg.PackageID, actualPkg.PackageID)
		assert.Equal(t, pkg.Sender, actualPkg.Sender)
		assert.Equal(t, pkg.Recipient, actualPkg.Recipient)
		assert.Equal(t, pkg.Origin, actualPkg.Origin)
		assert.Equal(t, pkg.Destination, actualPkg.Destination)
		assert.Equal(t, pkg.CurrentStatus, actualPkg.CurrentStatus)
	})

	t.Run("not found", func(t *testing.T) {
		pkg := &domain.Package{
			PackageID: "456",
		}

		mockService.On("UpdatePackage", pkg).Return(service.ErrEmptyPackageID)

		jsonData, _ := json.Marshal(pkg)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/packages/456", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, service.ErrEmptyPackageID.Error(), response.Error)
	})
}

func TestPackageHandler_DeletePackage(t *testing.T) {
	mockService := new(MockPackageService)
	handler := NewPackageHandler(mockService)
	router := setupTestRouter(handler)

	t.Run("success", func(t *testing.T) {
		mockService.On("DeletePackage", "123").Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/packages/123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("DeletePackage", "456").Return(service.ErrEmptyPackageID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/packages/456", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, service.ErrEmptyPackageID.Error(), response.Error)
	})
}
