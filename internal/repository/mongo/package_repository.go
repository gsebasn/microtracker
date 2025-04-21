package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/snavarro/microtracker/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrPackageNotFound = errors.New("package not found")
	defaultTimeout     = 5 * time.Second
)

type PackageRepository struct {
	collection *mongo.Collection
}

func NewPackageRepository(db *mongo.Database) *PackageRepository {
	return &PackageRepository{
		collection: db.Collection("packages"),
	}
}

func (r *PackageRepository) FindByID(id string) (*domain.Package, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var pkg domain.Package
	err := r.collection.FindOne(ctx, bson.M{"packageId": id}).Decode(&pkg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) FindAll(page, size int) ([]domain.Package, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	skip := int64((page - 1) * size)
	limit := int64(size)

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var packages []domain.Package
	if err = cursor.All(ctx, &packages); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *PackageRepository) Search(query string, page, size int) ([]domain.Package, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"$or": []bson.M{
			{"packageId": primitive.Regex{Pattern: query, Options: "i"}},
			{"sender.name": primitive.Regex{Pattern: query, Options: "i"}},
			{"recipient.name": primitive.Regex{Pattern: query, Options: "i"}},
			{"origin": primitive.Regex{Pattern: query, Options: "i"}},
			{"destination": primitive.Regex{Pattern: query, Options: "i"}},
			{"currentStatus": primitive.Regex{Pattern: query, Options: "i"}},
		},
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var packages []domain.Package
	if err = cursor.All(ctx, &packages); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *PackageRepository) Create(pkg *domain.Package) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pkg.CreatedAt = time.Now()
	pkg.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, pkg)
	return err
}

func (r *PackageRepository) Update(pkg *domain.Package) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pkg.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"packageId": pkg.PackageID},
		pkg,
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrPackageNotFound
	}

	return nil
}

func (r *PackageRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"packageId": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrPackageNotFound
	}

	return nil
}
