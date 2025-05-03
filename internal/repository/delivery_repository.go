package repository

import (
	"context"
	"errors"
	"fmt"
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

// GenerateRandomDeliveryPoints creates 40 delivery points with random distances
func (r *DeliveryPointRepository) GenerateRandomDeliveryPoints(ctx context.Context) error {
	// Check if delivery points already exist
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Delivery points already exist
	}

	// Generate 40 delivery points with IDs from 1001 to 1040
	points := make([]interface{}, 40)
	categories := []string{"Blue", "Green", "Yellow", "Purple"}

	// Create a base distance matrix for consistent distances
	distanceMatrix := make(map[int]map[int]int)

	// First, generate random distances between all points
	for i := 1001; i <= 1040; i++ {
		distanceMatrix[i] = make(map[int]int)
		for j := 1001; j <= 1040; j++ {
			if i != j { // No distance to itself
				// Generate a random distance between 25 and 100 km
				// Using a formula that ensures the distance is symmetric (same in both directions)
				// We use the sum of IDs as a seed to ensure symmetry
				seed := i + j
				distance := 25 + (seed % 76) // 25 + random value between 0 and 75
				distanceMatrix[i][j] = distance
			}
		}
	}

	for i := 0; i < 40; i++ {
		id := 1001 + i
		category := categories[i%4] // Distribute categories evenly

		// Create a delivery point with random distances to other points
		point := &models.DeliveryPoint{
			ID:          id,
			Category:    category,
			Pallets:     0, // Will be set by user
			Distances:   distanceMatrix[id],
			Description: fmt.Sprintf("Delivery Point %d", id),
			Address:     fmt.Sprintf("Address %d, Street %d", id-1000, (id-1000)%10),
		}

		points[i] = point
	}

	// Insert all delivery points in a single batch operation
	_, err = r.collection.InsertMany(ctx, points)
	return err
}
