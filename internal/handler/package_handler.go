package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snavarro/microtracker/internal/domain"
	"github.com/snavarro/microtracker/internal/repository/mongo"
	"github.com/snavarro/microtracker/internal/service"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Size    int         `json:"size,omitempty"`
	Success bool        `json:"success"`
}

type PackageHandler struct {
	service *service.PackageService
}

func NewPackageHandler(service *service.PackageService) *PackageHandler {
	return &PackageHandler{
		service: service,
	}
}

// @Summary Get a package by ID
// @Description Get package details by package ID
// @Tags packages
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /packages/{id} [get]
func (h *PackageHandler) GetPackage(c *gin.Context) {
	id := c.Param("id")
	pkg, err := h.service.GetPackage(id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrEmptyPackageID) {
			status = http.StatusBadRequest
		} else if errors.Is(err, mongo.ErrPackageNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, response{Error: err.Error(), Success: false})
		return
	}
	c.JSON(http.StatusOK, response{Data: pkg, Success: true})
}

// @Summary List all packages
// @Description Get a paginated list of all packages
// @Tags packages
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response
// @Failure 500 {object} response
// @Router /packages [get]
func (h *PackageHandler) ListPackages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	packages, total, err := h.service.ListPackages(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error(), Success: false})
		return
	}

	c.JSON(http.StatusOK, response{
		Data:    packages,
		Total:   total,
		Page:    page,
		Size:    size,
		Success: true,
	})
}

// @Summary Search packages
// @Description Search packages with pagination
// @Tags packages
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response
// @Failure 500 {object} response
// @Router /packages/search [get]
func (h *PackageHandler) SearchPackages(c *gin.Context) {
	query := c.Query("query")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	packages, total, err := h.service.SearchPackages(query, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error(), Success: false})
		return
	}

	c.JSON(http.StatusOK, response{
		Data:    packages,
		Total:   total,
		Page:    page,
		Size:    size,
		Success: true,
	})
}

// @Summary Create a new package
// @Description Create a new package
// @Tags packages
// @Accept json
// @Produce json
// @Param package body domain.Package true "Package details"
// @Success 201 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /packages [post]
func (h *PackageHandler) CreatePackage(c *gin.Context) {
	var pkg domain.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "Invalid request body", Success: false})
		return
	}

	if err := h.service.CreatePackage(&pkg); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidPackage) || errors.Is(err, service.ErrEmptyPackageID) {
			status = http.StatusBadRequest
		}
		c.JSON(status, response{Error: err.Error(), Success: false})
		return
	}

	c.JSON(http.StatusCreated, response{Data: pkg, Success: true})
}

// @Summary Update a package
// @Description Update an existing package
// @Tags packages
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Param package body domain.Package true "Package details"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /packages/{id} [put]
func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	id := c.Param("id")
	var pkg domain.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, response{Error: "Invalid request body", Success: false})
		return
	}

	pkg.PackageID = id
	if err := h.service.UpdatePackage(&pkg); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidPackage) || errors.Is(err, service.ErrEmptyPackageID) {
			status = http.StatusBadRequest
		} else if errors.Is(err, mongo.ErrPackageNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, response{Error: err.Error(), Success: false})
		return
	}

	c.JSON(http.StatusOK, response{Data: pkg, Success: true})
}

// @Summary Delete a package
// @Description Delete a package by ID
// @Tags packages
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Success 204
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /packages/{id} [delete]
func (h *PackageHandler) DeletePackage(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePackage(id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrEmptyPackageID) {
			status = http.StatusBadRequest
		} else if errors.Is(err, mongo.ErrPackageNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, response{Error: err.Error(), Success: false})
		return
	}
	c.Status(http.StatusNoContent)
}
