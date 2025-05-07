package repository

import (
	"context"
	"errors"
	"github.com/vladislavprovich/knuca-diploma-work/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RouteRepository handles database operations for delivery routes
type RouteRepository struct {
	collection *mongo.Collection
}

// NewRouteRepository creates a new repository for routes
func NewRouteRepository(db *mongo.Database) *RouteRepository {
	return &RouteRepository{
		collection: db.Collection("routes"),
	}
}

// Create adds a new route to the database
func (r *RouteRepository) Create(ctx context.Context, route *models.Route) error {
	// Set creation and update timestamps
	now := time.Now()
	route.CreatedAt = now
	route.UpdatedAt = now

	// Check if a route with this ID already exists
	var existingRoute models.Route
	err := r.collection.FindOne(ctx, bson.M{"_id": route.ID}).Decode(&existingRoute)
	if err == nil {
		// Route with this ID already exists, update it instead
		_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": route.ID}, route)
		return err
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		// An error occurred that wasn't just "no documents found"
		return err
	}

	// No existing route with this ID, create a new one
	_, err = r.collection.InsertOne(ctx, route)
	return err
}

// GetByID retrieves a route by its ID
func (r *RouteRepository) GetByID(ctx context.Context, id int) (*models.Route, error) {
	var route models.Route
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&route)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("route not found")
		}
		return nil, err
	}
	return &route, nil
}

// GetAll retrieves all routes
func (r *RouteRepository) GetAll(ctx context.Context) ([]*models.Route, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

// GetByStatus retrieves routes by their status
func (r *RouteRepository) GetByStatus(ctx context.Context, status string) ([]*models.Route, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

// GetByDate retrieves routes for a specific date
func (r *RouteRepository) GetByDate(ctx context.Context, date time.Time) ([]*models.Route, error) {
	// Create start and end of the day for the query
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	cursor, err := r.collection.Find(ctx, bson.M{
		"delivery_date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

// GetByTruckID retrieves routes assigned to a specific truck
func (r *RouteRepository) GetByTruckID(ctx context.Context, truckID int) ([]*models.Route, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"truck_id": truckID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

// Update updates an existing route
func (r *RouteRepository) Update(ctx context.Context, route *models.Route) error {
	// Update the timestamp
	route.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": route.ID}, route)
	return err
}

// UpdateStatus updates only the status of a route
func (r *RouteRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		}},
	)
	return err
}

// Delete removes a route from the database
func (r *RouteRepository) Delete(ctx context.Context, id int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteAll removes all routes from the database
func (r *RouteRepository) DeleteAll(ctx context.Context) error {
	// First, get all routes to find associated trucks
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var routes []*models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return err
	}

	// Delete all routes
	_, err = r.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

// GetNextRouteID gets the next available route ID
func (r *RouteRepository) GetNextRouteID(ctx context.Context) (int, error) {
	// Find the route with the highest ID
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	var route models.Route

	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&route)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 1, nil // Start with ID 1 if no routes exist
		}
		return 0, err
	}

	return route.ID + 1, nil
}
