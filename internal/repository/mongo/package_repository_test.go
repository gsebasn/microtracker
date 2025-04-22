package mongo

import (
	"testing"
	"time"

	"github.com/snavarro/microtracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestPackageRepository_FindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		expectedPkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "packageId", Value: expectedPkg.PackageID},
			{Key: "sender", Value: bson.D{
				{Key: "name", Value: expectedPkg.Sender.Name},
				{Key: "address", Value: expectedPkg.Sender.Address},
			}},
			{Key: "createdAt", Value: expectedPkg.CreatedAt},
			{Key: "updatedAt", Value: expectedPkg.UpdatedAt},
		}))

		pkg, err := repo.FindByID("123")
		require.NoError(t, err)
		assert.Equal(t, expectedPkg.PackageID, pkg.PackageID)
		assert.Equal(t, expectedPkg.Sender, pkg.Sender)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		pkg, err := repo.FindByID("456")
		assert.Error(t, err)
		assert.Nil(t, pkg)
	})
}

func TestPackageRepository_FindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
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

		mt.AddMockResponses(
			mtest.CreateCursorResponse(2, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "packageId", Value: expectedPackages[0].PackageID},
				{Key: "sender", Value: bson.D{
					{Key: "name", Value: expectedPackages[0].Sender.Name},
					{Key: "address", Value: expectedPackages[0].Sender.Address},
				}},
			}, bson.D{
				{Key: "packageId", Value: expectedPackages[1].PackageID},
				{Key: "sender", Value: bson.D{
					{Key: "name", Value: expectedPackages[1].Sender.Name},
					{Key: "address", Value: expectedPackages[1].Sender.Address},
				}},
			}),
			mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch),
		)

		packages, total, err := repo.FindAll(1, 10)
		require.NoError(t, err)
		assert.Equal(t, expectedPackages, packages)
		assert.Equal(t, int64(2), total)
	})

	mt.Run("error", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		packages, total, err := repo.FindAll(1, 10)
		assert.Error(t, err)
		assert.Empty(t, packages)
		assert.Equal(t, int64(0), total)
	})
}

func TestPackageRepository_Search(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		expectedPackages := []domain.Package{
			{
				PackageID: "123",
				Sender: domain.Address{
					Name:    "John Doe",
					Address: "123 Main St",
				},
			},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
				{Key: "packageId", Value: expectedPackages[0].PackageID},
				{Key: "sender", Value: bson.D{
					{Key: "name", Value: expectedPackages[0].Sender.Name},
					{Key: "address", Value: expectedPackages[0].Sender.Address},
				}},
			}),
			mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch),
		)

		packages, total, err := repo.Search("John", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, expectedPackages, packages)
		assert.Equal(t, int64(1), total)
	})

	mt.Run("error", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		packages, total, err := repo.Search("test", 1, 10)
		assert.Error(t, err)
		assert.Empty(t, packages)
		assert.Equal(t, int64(0), total)
	})
}

func TestPackageRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		pkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}})

		err := repo.Create(pkg)
		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		pkg := &domain.Package{
			PackageID: "123",
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		err := repo.Create(pkg)
		assert.Error(t, err)
	})
}

func TestPackageRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		pkg := &domain.Package{
			PackageID: "123",
			Sender: domain.Address{
				Name:    "John Doe",
				Address: "123 Main St",
			},
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}})

		err := repo.Update(pkg)
		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		pkg := &domain.Package{
			PackageID: "123",
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		err := repo.Update(pkg)
		assert.Error(t, err)
	})
}

func TestPackageRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}})

		err := repo.Delete("123")
		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		repo := NewPackageRepository(mt.DB)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		err := repo.Delete("123")
		assert.Error(t, err)
	})
}
