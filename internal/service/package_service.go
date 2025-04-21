package service

import (
	"errors"
	"strings"

	"github.com/snavarro/microtracker/internal/domain"
)

var (
	ErrInvalidPackage = errors.New("invalid package data")
	ErrEmptyPackageID = errors.New("package ID cannot be empty")
)

type PackageService struct {
	repo domain.PackageRepository
}

func NewPackageService(repo domain.PackageRepository) *PackageService {
	return &PackageService{
		repo: repo,
	}
}

func (s *PackageService) GetPackage(id string) (*domain.Package, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrEmptyPackageID
	}
	return s.repo.FindByID(id)
}

func (s *PackageService) ListPackages(page, size int) ([]domain.Package, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return s.repo.FindAll(page, size)
}

func (s *PackageService) SearchPackages(query string, page, size int) ([]domain.Package, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return s.repo.Search(query, page, size)
}

func (s *PackageService) CreatePackage(pkg *domain.Package) error {
	if err := validatePackage(pkg); err != nil {
		return err
	}
	return s.repo.Create(pkg)
}

func (s *PackageService) UpdatePackage(pkg *domain.Package) error {
	if err := validatePackage(pkg); err != nil {
		return err
	}
	return s.repo.Update(pkg)
}

func (s *PackageService) DeletePackage(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrEmptyPackageID
	}
	return s.repo.Delete(id)
}

func validatePackage(pkg *domain.Package) error {
	if pkg == nil {
		return ErrInvalidPackage
	}

	if strings.TrimSpace(pkg.PackageID) == "" {
		return ErrEmptyPackageID
	}

	if strings.TrimSpace(pkg.Sender.Name) == "" || strings.TrimSpace(pkg.Sender.Address) == "" {
		return errors.New("sender name and address are required")
	}

	if strings.TrimSpace(pkg.Recipient.Name) == "" || strings.TrimSpace(pkg.Recipient.Address) == "" {
		return errors.New("recipient name and address are required")
	}

	if strings.TrimSpace(pkg.Origin) == "" {
		return errors.New("origin is required")
	}

	if strings.TrimSpace(pkg.Destination) == "" {
		return errors.New("destination is required")
	}

	if strings.TrimSpace(pkg.CurrentStatus) == "" {
		return errors.New("current status is required")
	}

	return nil
}
