package domain

import "time"

type Address struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type Event struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Location  string    `json:"location" bson:"location"`
	Status    string    `json:"status" bson:"status"`
}

type Package struct {
	PackageID     string    `json:"packageId" bson:"packageId"`
	Sender        Address   `json:"sender" bson:"sender"`
	Recipient     Address   `json:"recipient" bson:"recipient"`
	Origin        string    `json:"origin" bson:"origin"`
	Destination   string    `json:"destination" bson:"destination"`
	CurrentStatus string    `json:"currentStatus" bson:"currentStatus"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
	Events        []Event   `json:"events,omitempty" bson:"events,omitempty"`
}

type PackageRepository interface {
	FindByID(id string) (*Package, error)
	FindAll(page, size int) ([]Package, int64, error)
	Search(query string, page, size int) ([]Package, int64, error)
	Create(pkg *Package) error
	Update(pkg *Package) error
	Delete(id string) error
}
