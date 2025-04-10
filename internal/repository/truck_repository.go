package repository

import (
	"context"
	"errors"
	"universati-savokh/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TruckRepository handles database operations for trucks
type TruckRepository struct {
	collection *mongo.Collection
}

// NewTruckRepository creates a new repository for trucks
func NewTruckRepository(db *mongo.Database) *TruckRepository {
	return &TruckRepository{
		collection: db.Collection("trucks"),
	}
}

// Create adds a new truck to the database
func (r *TruckRepository) Create(ctx context.Context, truck *models.Truck) error {
	_, err := r.collection.InsertOne(ctx, truck)
	return err
}

// GetByID retrieves a truck by its ID
func (r *TruckRepository) GetByID(ctx context.Context, id int) (*models.Truck, error) {
	var truck models.Truck
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&truck)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("truck not found")
		}
		return nil, err
	}
	return &truck, nil
}

// GetAll retrieves all trucks
func (r *TruckRepository) GetAll(ctx context.Context) ([]*models.Truck, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trucks []*models.Truck
	if err = cursor.All(ctx, &trucks); err != nil {
		return nil, err
	}

	return trucks, nil
}

// GetByCapacity retrieves trucks by their capacity
func (r *TruckRepository) GetByCapacity(ctx context.Context, capacity int) ([]*models.Truck, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"capacity": capacity})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trucks []*models.Truck
	if err = cursor.All(ctx, &trucks); err != nil {
		return nil, err
	}

	return trucks, nil
}

// Update updates an existing truck
func (r *TruckRepository) Update(ctx context.Context, truck *models.Truck) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": truck.ID}, truck)
	return err
}

// GetAvailableTrucks retrieves all available trucks
func (r *TruckRepository) GetAvailableTrucks(ctx context.Context) ([]*models.Truck, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"available": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trucks []*models.Truck
	if err = cursor.All(ctx, &trucks); err != nil {
		return nil, err
	}

	return trucks, nil
}

// GetNextTruckID gets the next available truck ID
func (r *TruckRepository) GetNextTruckID(ctx context.Context) (int, error) {
	// Find the truck with the highest ID
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	var truck models.Truck

	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&truck)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 1, nil // Start with ID 1 if no trucks exist
		}
		return 0, err
	}

	return truck.ID + 1, nil
}

// InitializeFleet creates a default fleet of trucks if none exist
func (r *TruckRepository) InitializeFleet(ctx context.Context) error {
	// Check if trucks already exist
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Trucks already exist
	}

	// Create default fleet with different capacities
	trucks := []interface{}{
		// 33-pallet trucks (with trailers)
		&models.Truck{
			ID:           1,
			LicensePlate: "ABC-123",
			Brand:        "Volvo",
			Capacity:     33,
			DriverName:   "John",
			DriverSurname: "Doe",
			HasTrailer:   true,
			TrailerInfo:  "Trailer-1",
			Available:    true,
		},
		&models.Truck{
			ID:           2,
			LicensePlate: "DEF-456",
			Brand:        "Scania",
			Capacity:     33,
			DriverName:   "Jane",
			DriverSurname: "Smith",
			HasTrailer:   true,
			TrailerInfo:  "Trailer-2",
			Available:    true,
		},
		
		// 18-pallet trucks
		&models.Truck{
			ID:           3,
			LicensePlate: "GHI-789",
			Brand:        "Mercedes",
			Capacity:     18,
			DriverName:   "Michael",
			DriverSurname: "Johnson",
			HasTrailer:   false,
			Available:    true,
		},
		&models.Truck{
			ID:           4,
			LicensePlate: "JKL-012",
			Brand:        "MAN",
			Capacity:     18,
			DriverName:   "Emily",
			DriverSurname: "Brown",
			HasTrailer:   false,
			Available:    true,
		},
		
		// 15-pallet trucks
		&models.Truck{
			ID:           5,
			LicensePlate: "MNO-345",
			Brand:        "Iveco",
			Capacity:     15,
			DriverName:   "David",
			DriverSurname: "Wilson",
			HasTrailer:   false,
			Available:    true,
		},
		&models.Truck{
			ID:           6,
			LicensePlate: "PQR-678",
			Brand:        "DAF",
			Capacity:     15,
			DriverName:   "Sarah",
			DriverSurname: "Taylor",
			HasTrailer:   false,
			Available:    true,
		},
		
		// 10-pallet trucks
		&models.Truck{
			ID:           7,
			LicensePlate: "STU-901",
			Brand:        "Renault",
			Capacity:     10,
			DriverName:   "Robert",
			DriverSurname: "Anderson",
			HasTrailer:   false,
			Available:    true,
		},
		&models.Truck{
			ID:           8,
			LicensePlate: "VWX-234",
			Brand:        "Fiat",
			Capacity:     10,
			DriverName:   "Lisa",
			DriverSurname: "Martinez",
			HasTrailer:   false,
			Available:    true,
		},
	}

	// Insert all trucks in a single batch operation
	_, err = r.collection.InsertMany(ctx, trucks)
	return err
}
