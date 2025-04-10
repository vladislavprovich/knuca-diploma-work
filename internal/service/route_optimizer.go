package service

import (
	"context"
	"errors"
	"github.com/vladislavprovich/knuca-diploma-work/internal/models"
	"github.com/vladislavprovich/knuca-diploma-work/internal/repository"
	"sort"
	"time"
)

// RouteOptimizerService handles route optimization logic
type RouteOptimizerService struct {
	DeliveryRepo *repository.DeliveryPointRepository
	TruckRepo    *repository.TruckRepository
	RouteRepo    *repository.RouteRepository
}

// NewRouteOptimizerService creates a new route optimizer service
func NewRouteOptimizerService(
	deliveryRepo *repository.DeliveryPointRepository,
	truckRepo *repository.TruckRepository,
	routeRepo *repository.RouteRepository,
) *RouteOptimizerService {
	return &RouteOptimizerService{
		DeliveryRepo: deliveryRepo,
		TruckRepo:    truckRepo,
		RouteRepo:    routeRepo,
	}
}

// OptimizeRoutes generates optimized delivery routes based on delivery points with pallets
func (s *RouteOptimizerService) OptimizeRoutes(ctx context.Context, deliveryDate time.Time) ([]*models.Route, error) {
	// Get all delivery points
	deliveryPoints, err := s.DeliveryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Filter out delivery points with no pallets
	pointsWithPallets := make([]*models.DeliveryPoint, 0)
	for _, point := range deliveryPoints {
		if point.Pallets > 0 {
			pointsWithPallets = append(pointsWithPallets, point)
		}
	}

	if len(pointsWithPallets) == 0 {
		return nil, errors.New("no delivery points with pallets found")
	}

	// Get all available trucks
	availableTrucks, err := s.TruckRepo.GetAvailableTrucks(ctx)
	if err != nil {
		return nil, err
	}

	if len(availableTrucks) == 0 {
		return nil, errors.New("no available trucks found")
	}

	// Sort trucks by capacity (largest first)
	sort.Slice(availableTrucks, func(i, j int) bool {
		return availableTrucks[i].Capacity > availableTrucks[j].Capacity
	})

	// Create a map of delivery points for easier access
	pointMap := make(map[int]*models.DeliveryPoint)
	for _, point := range deliveryPoints {
		pointMap[point.ID] = point
	}

	// Group delivery points by category
	categoryPoints := make(map[string][]*models.DeliveryPoint)
	for _, point := range pointsWithPallets {
		categoryPoints[point.Category] = append(categoryPoints[point.Category], point)
	}

	// Generate routes
	routes := make([]*models.Route, 0)

	// Start with Blue category (33-pallet trucks)
	if bluePoints, ok := categoryPoints["Blue"]; ok && len(bluePoints) > 0 {
		blueRoutes, err := s.generateRoutesForCategory(ctx, bluePoints, availableTrucks, 33, pointMap, deliveryDate)
		if err != nil {
			return nil, err
		}
		routes = append(routes, blueRoutes...)
	}

	// Then Green category (18-pallet trucks)
	if greenPoints, ok := categoryPoints["Green"]; ok && len(greenPoints) > 0 {
		greenRoutes, err := s.generateRoutesForCategory(ctx, greenPoints, availableTrucks, 18, pointMap, deliveryDate)
		if err != nil {
			return nil, err
		}
		routes = append(routes, greenRoutes...)
	}

	// Then Yellow category (15-pallet trucks)
	if yellowPoints, ok := categoryPoints["Yellow"]; ok && len(yellowPoints) > 0 {
		yellowRoutes, err := s.generateRoutesForCategory(ctx, yellowPoints, availableTrucks, 15, pointMap, deliveryDate)
		if err != nil {
			return nil, err
		}
		routes = append(routes, yellowRoutes...)
	}

	// Finally Purple category (10-pallet trucks)
	if purplePoints, ok := categoryPoints["Purple"]; ok && len(purplePoints) > 0 {
		purpleRoutes, err := s.generateRoutesForCategory(ctx, purplePoints, availableTrucks, 10, pointMap, deliveryDate)
		if err != nil {
			return nil, err
		}
		routes = append(routes, purpleRoutes...)
	}

	return routes, nil
}

// generateRoutesForCategory creates routes for delivery points of a specific category
func (s *RouteOptimizerService) generateRoutesForCategory(
	ctx context.Context,
	points []*models.DeliveryPoint,
	availableTrucks []*models.Truck,
	requiredCapacity int,
	pointMap map[int]*models.DeliveryPoint,
	deliveryDate time.Time,
) ([]*models.Route, error) {
	routes := make([]*models.Route, 0)

	// Find trucks with the required capacity
	suitableTrucks := make([]*models.Truck, 0)
	for _, truck := range availableTrucks {
		if truck.Capacity == requiredCapacity && truck.Available {
			suitableTrucks = append(suitableTrucks, truck)
		}
	}

	// If no trucks with exact capacity, try smaller trucks for Yellow category
	if len(suitableTrucks) == 0 && requiredCapacity == 15 {
		for _, truck := range availableTrucks {
			if truck.Capacity == 10 && truck.Available {
				suitableTrucks = append(suitableTrucks, truck)
			}
		}
	}

	if len(suitableTrucks) == 0 {
		return nil, errors.New("no suitable trucks available for this category")
	}

	// Sort points by number of pallets (descending)
	sort.Slice(points, func(i, j int) bool {
		return points[i].Pallets > points[j].Pallets
	})

	// Create routes using nearest neighbor algorithm
	remaining := make([]*models.DeliveryPoint, len(points))
	copy(remaining, points)

	truckIndex := 0
	for len(remaining) > 0 && truckIndex < len(suitableTrucks) {
		truck := suitableTrucks[truckIndex]

		// Get next route ID
		routeID, err := s.RouteRepo.GetNextRouteID(ctx)
		if err != nil {
			return nil, err
		}

		// Start with the first point
		startPoint := remaining[0]
		remaining = remaining[1:]

		// Create a new route
		route := &models.Route{
			ID:             routeID,
			TruckID:        truck.ID,
			DeliveryDate:   deliveryDate,
			StartPoint:     startPoint.ID,
			DeliveryPoints: []int{startPoint.ID},
			Status:         models.RouteStatus.Planned,
		}

		// Mark truck as unavailable
		truck.Available = false

		// Add more points using nearest neighbor algorithm
		currentPoint := startPoint
		for len(remaining) > 0 {
			// Find the nearest point
			nearestIndex := 0
			nearestDistance := -1

			for i, point := range remaining {
				if distance, exists := currentPoint.Distances[point.ID]; exists {
					if nearestDistance == -1 || distance < nearestDistance {
						nearestDistance = distance
						nearestIndex = i
					}
				}
			}

			// Add the nearest point to the route
			nextPoint := remaining[nearestIndex]
			route.DeliveryPoints = append(route.DeliveryPoints, nextPoint.ID)

			// Remove the point from remaining
			remaining = append(remaining[:nearestIndex], remaining[nearestIndex+1:]...)

			// Update current point
			currentPoint = nextPoint

			// Break if we've reached a reasonable number of points per route
			if len(route.DeliveryPoints) >= 5 {
				break
			}
		}

		// Calculate total distance and cost
		route.CalculateTotalDistance(pointMap)
		route.CalculateTotalCost(truck)

		// Add route to the list
		routes = append(routes, route)

		// Move to the next truck
		truckIndex++
	}

	return routes, nil
}

// SaveRoutes saves the generated routes to the database
func (s *RouteOptimizerService) SaveRoutes(ctx context.Context, routes []*models.Route) error {
	for _, route := range routes {
		if err := s.RouteRepo.Create(ctx, route); err != nil {
			return err
		}

		// Update truck availability
		truck, err := s.TruckRepo.GetByID(ctx, route.TruckID)
		if err != nil {
			return err
		}

		truck.Available = false
		if err := s.TruckRepo.Update(ctx, truck); err != nil {
			return err
		}
	}

	return nil
}

// UpdateRoute updates an existing route
func (s *RouteOptimizerService) UpdateRoute(ctx context.Context, route *models.Route) error {
	// Get the truck for this route
	truck, err := s.TruckRepo.GetByID(ctx, route.TruckID)
	if err != nil {
		return err
	}

	// Get all delivery points
	deliveryPoints, err := s.DeliveryRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Create a map of delivery points for easier access
	pointMap := make(map[int]*models.DeliveryPoint)
	for _, point := range deliveryPoints {
		pointMap[point.ID] = point
	}

	// Validate the route
	if !route.ValidateRoute(truck, pointMap) {
		return errors.New("invalid route: truck capacity does not match delivery point requirements")
	}

	// Recalculate distance and cost
	route.CalculateTotalDistance(pointMap)
	route.CalculateTotalCost(truck)

	// Update the route
	return s.RouteRepo.Update(ctx, route)
}
