package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snavarro/microtracker/config"
	"github.com/snavarro/microtracker/internal/handler"
	"github.com/snavarro/microtracker/internal/middleware"
	mongorepo "github.com/snavarro/microtracker/internal/repository/mongo"
	"github.com/snavarro/microtracker/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// Rate limits from configuration
	listEndpointLimit   = 200
	searchEndpointLimit = 150
	createEndpointLimit = 50
	rateLimitTTL        = 5 * time.Minute

	// Test endpoints
	listEndpoint   = "/api/v1/packages"
	searchEndpoint = "/api/v1/packages/search?query=test"
)

func setupTestServer(t *testing.T) (*gin.Engine, *mongo.Database) {
	// Set test environment
	os.Setenv("APP_ENV", "test")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load test configuration
	cfg, err := config.NewConfig()
	require.NoError(t, err, "Failed to load configuration")

	// Connect to test database
	db, err := config.ConnectDB(cfg)
	require.NoError(t, err, "Failed to connect to database")

	// Initialize components
	packageRepo := mongorepo.NewPackageRepository(db)
	packageService := service.NewPackageService(packageRepo)
	packageHandler := handler.NewPackageHandler(packageService)

	// Create router
	router := gin.New()
	router.Use(gin.Recovery())

	// Add rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(&cfg.RateLimit)
	router.Use(rateLimiter.RateLimit())

	// Setup routes
	api := router.Group("/api/v1")
	{
		packages := api.Group("/packages")
		{
			packages.GET("", packageHandler.ListPackages)
			packages.GET("/search", packageHandler.SearchPackages)
			packages.GET("/:id", packageHandler.GetPackage)
			packages.POST("", packageHandler.CreatePackage)
			packages.PUT("/:id", packageHandler.UpdatePackage)
			packages.DELETE("/:id", packageHandler.DeletePackage)
		}
	}

	return router, db
}

func cleanupDatabase(t *testing.T, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Drop(ctx)
	require.NoError(t, err, "Failed to cleanup test database")

	err = db.Client().Disconnect(ctx)
	require.NoError(t, err, "Failed to disconnect from database")
}

func TestRateLimitListPackages(t *testing.T) {
	router, db := setupTestServer(t)
	defer cleanupDatabase(t, db)

	// Make requests up to the limit
	for i := 0; i < listEndpointLimit+10; i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, listEndpoint, nil)
		require.NoError(t, err)

		router.ServeHTTP(w, req)

		if i < listEndpointLimit {
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited", i)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "Rate limit exceeded", response["error"])
			assert.Equal(t, false, response["success"])
			break
		}
	}
}

func TestRateLimitSearchPackages(t *testing.T) {
	router, db := setupTestServer(t)
	defer cleanupDatabase(t, db)

	for i := 0; i < searchEndpointLimit+10; i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, searchEndpoint, nil)
		require.NoError(t, err)

		router.ServeHTTP(w, req)

		if i < searchEndpointLimit {
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited", i)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "Rate limit exceeded", response["error"])
			assert.Equal(t, false, response["success"])
			break
		}
	}
}

func TestRateLimitCreatePackage(t *testing.T) {
	router, db := setupTestServer(t)
	defer cleanupDatabase(t, db)

	packageData := map[string]interface{}{
		"packageId": "TEST123",
		"sender": map[string]string{
			"name":    "Test Sender",
			"address": "123 Test St",
		},
		"recipient": map[string]string{
			"name":    "Test Recipient",
			"address": "456 Test St",
		},
		"origin":        "Test Origin",
		"destination":   "Test Destination",
		"currentStatus": "created",
	}

	for i := 0; i < createEndpointLimit+10; i++ {
		jsonData, err := json.Marshal(packageData)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, listEndpoint, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		if i < createEndpointLimit {
			assert.Equal(t, http.StatusCreated, w.Code, "Request %d should succeed", i)
			// Modify package ID to ensure unique entries
			packageData["packageId"] = fmt.Sprintf("TEST%d", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited", i)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "Rate limit exceeded", response["error"])
			assert.Equal(t, false, response["success"])
			break
		}
	}
}

func TestRateLimitReset(t *testing.T) {
	router, db := setupTestServer(t)
	defer cleanupDatabase(t, db)

	// Make requests up to the limit
	for i := 0; i < listEndpointLimit+1; i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, listEndpoint, nil)
		require.NoError(t, err)
		router.ServeHTTP(w, req)

		if i < listEndpointLimit {
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited", i)
			break
		}
	}

	// Wait for rate limit to reset (TTL + 1 second buffer)
	time.Sleep(rateLimitTTL + time.Second)

	// Verify rate limit has reset
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, listEndpoint, nil)
	require.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Request should succeed after rate limit reset")
}

func TestRateLimitDifferentIPs(t *testing.T) {
	router, db := setupTestServer(t)
	defer cleanupDatabase(t, db)

	// Test IPs
	ip1 := "192.168.1.1:12345"
	ip2 := "192.168.1.2:12345"

	// Test first IP until rate limited
	for i := 0; i < listEndpointLimit+1; i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, listEndpoint, nil)
		require.NoError(t, err)
		req.RemoteAddr = ip1
		router.ServeHTTP(w, req)

		if i < listEndpointLimit {
			assert.Equal(t, http.StatusOK, w.Code, "Request %d from IP1 should succeed", i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d from IP1 should be rate limited", i)
			break
		}
	}

	// Test second IP - should work regardless of first IP being rate limited
	for i := 0; i < 5; i++ { // Test multiple requests from second IP
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, listEndpoint, nil)
		require.NoError(t, err)
		req.RemoteAddr = ip2
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d from IP2 should succeed", i)
	}
}
