package repository

import (
	"context"
	"errors"

	"github.com/vladislavprovich/knuca-diploma-work/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeliveryPointRepository handles database operations for delivery points
type DeliveryPointRepository struct {
	collection *mongo.Collection
}

// NewDeliveryPointRepository creates a new repository for delivery points
func NewDeliveryPointRepository(db *mongo.Database) *DeliveryPointRepository {
	return &DeliveryPointRepository{
		collection: db.Collection("delivery_points"),
	}
}

// Create adds a new delivery point to the database
func (r *DeliveryPointRepository) Create(ctx context.Context, point *models.DeliveryPoint) error {
	_, err := r.collection.InsertOne(ctx, point)
	return err
}

// GetByID retrieves a delivery point by its ID
func (r *DeliveryPointRepository) GetByID(ctx context.Context, id int) (*models.DeliveryPoint, error) {
	var point models.DeliveryPoint
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&point)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("delivery point not found")
		}
		return nil, err
	}
	return &point, nil
}

// GetAll retrieves all delivery points
func (r *DeliveryPointRepository) GetAll(ctx context.Context) ([]*models.DeliveryPoint, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var points []*models.DeliveryPoint
	if err = cursor.All(ctx, &points); err != nil {
		return nil, err
	}

	return points, nil
}

// GetByCategory retrieves delivery points by category
func (r *DeliveryPointRepository) GetByCategory(ctx context.Context, category string) ([]*models.DeliveryPoint, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var points []*models.DeliveryPoint
	if err = cursor.All(ctx, &points); err != nil {
		return nil, err
	}

	return points, nil
}

// Update updates an existing delivery point
func (r *DeliveryPointRepository) Update(ctx context.Context, point *models.DeliveryPoint) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": point.ID}, point)
	return err
}

// Delete removes a delivery point from the database
func (r *DeliveryPointRepository) Delete(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GenerateRandomDeliveryPoints creates delivery points with real Kyiv locations and accurate distances
func (r *DeliveryPointRepository) GenerateRandomDeliveryPoints(ctx context.Context) error {
	// Check if delivery points already exist
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Delivery points already exist
	}

	// Get real Kyiv locations
	kyivLocations := GetKyivLocations()

	// Calculate real distances between locations
	distanceMatrix := CalculateRealDistances(kyivLocations)

	// Create points from Kyiv locations
	points := make([]interface{}, len(kyivLocations))

	for i, location := range kyivLocations {
		id := 1001 + i

		point := &models.DeliveryPoint{
			ID:          id,
			Category:    location.Category,
			Pallets:     0, // Will be set by user
			Distances:   distanceMatrix[id],
			Description: location.Name,
			Address:     location.Address,
			Latitude:    location.Latitude,
			Longitude:   location.Longitude,
		}

		points[i] = point
	}

	// Insert all delivery points in a single batch operation
	_, err = r.collection.InsertMany(ctx, points)
	return err
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
