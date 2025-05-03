package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/vladislavprovich/knuca-diploma-work/internal/models"
	"github.com/vladislavprovich/knuca-diploma-work/internal/repository"
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
		return nil, fmt.Errorf("failed to get delivery points: %w", err)
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
		return nil, fmt.Errorf("failed to get available trucks: %w", err)
	}

	if len(availableTrucks) == 0 {
		// If no trucks are available, try to reset all trucks to available
		log.Println("No available trucks found, attempting to reset truck availability")
		allTrucks, err := s.TruckRepo.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get all trucks: %w", err)
		}

		// Make all trucks available
		for _, truck := range allTrucks {
			truck.Available = true
			if err := s.TruckRepo.Update(ctx, truck); err != nil {
				log.Printf("Warning: Failed to reset availability for truck %d: %v", truck.ID, err)
			}
		}

		// Get available trucks again
		availableTrucks, err = s.TruckRepo.GetAvailableTrucks(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get available trucks after reset: %w", err)
		}

		if len(availableTrucks) == 0 {
			return nil, errors.New("no trucks available in the system")
		}

		log.Printf("Successfully reset %d trucks to available", len(availableTrucks))
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

	// Create a map to track remaining pallets for each delivery point
	remainingPalletsMap := make(map[int]int)
	for _, point := range pointsWithPallets {
		remainingPalletsMap[point.ID] = point.Pallets
	}

	// Check for delivery points with pallets exceeding the largest truck capacity
	// and split them into multiple deliveries
	largestTruckCapacity := 0
	if len(availableTrucks) > 0 {
		largestTruckCapacity = availableTrucks[0].Capacity
	}

	// Create a map to track which points have been assigned to routes
	assignedPoints := make(map[int]bool)

	// Process delivery points by category, starting with Blue (which needs the largest trucks)
	categories := []string{"Blue", "Green", "Yellow", "Purple"}

	// First, create a map of points by category
	categoryPoints := make(map[string][]*models.DeliveryPoint)
	for _, point := range pointsWithPallets {
		if !assignedPoints[point.ID] {
			// Only add the point if it hasn't been fully processed yet
			if point.Pallets <= largestTruckCapacity || remainingPalletsMap[point.ID] > 0 {
				categoryPoints[point.Category] = append(categoryPoints[point.Category], point)
			}
		}
	}

	// Generate routes
	routes := make([]*models.Route, 0)

	// Process categories in order (Blue, Green, Yellow, Purple)
	for _, category := range categories {
		if points, ok := categoryPoints[category]; ok && len(points) > 0 {
			log.Printf("Processing %d %s category points", len(points), category)

			// Determine the preferred truck capacity for this category
			preferredCapacity := 0
			var alternativeCapacities []int

			switch category {
			case "Blue":
				preferredCapacity = 33
				alternativeCapacities = []int{18, 15}
			case "Green":
				preferredCapacity = 18
				alternativeCapacities = []int{15, 10}
			case "Yellow":
				preferredCapacity = 15
				alternativeCapacities = []int{10}
			case "Purple":
				preferredCapacity = 10
				alternativeCapacities = []int{}
			}

			// Try to use the preferred truck capacity first
			categoryRoutes, err := s.generateRoutesForCategory(ctx, points, availableTrucks, preferredCapacity, pointMap, deliveryDate)
			if err != nil {
				// If preferred trucks aren't available, try alternative capacities
				log.Printf("Could not use %d-pallet trucks for %s category: %v. Trying alternative trucks.",
					preferredCapacity, category, err)

				if len(alternativeCapacities) > 0 {
					splitRoutes, splitErr := s.generateSplitRoutesForCategory(ctx, points, availableTrucks,
						alternativeCapacities, pointMap, deliveryDate)
					if splitErr == nil {
						log.Printf("Successfully created %d routes with alternative trucks for %s category",
							len(splitRoutes), category)
						routes = append(routes, splitRoutes...)

						// Mark points as assigned
						for _, route := range splitRoutes {
							for _, pointID := range route.DeliveryPoints {
								assignedPoints[pointID] = true
							}
						}
					} else {
						log.Printf("Failed to create routes with alternative trucks for %s category: %v",
							category, splitErr)
					}
				} else {
					log.Printf("No alternative truck capacities available for %s category", category)
				}
			} else {
				log.Printf("Successfully created %d routes with %d-pallet trucks for %s category",
					len(categoryRoutes), preferredCapacity, category)
				routes = append(routes, categoryRoutes...)

				// Mark points as assigned
				for _, route := range categoryRoutes {
					for _, pointID := range route.DeliveryPoints {
						assignedPoints[pointID] = true
					}
				}
			}
		}
	}

	// Check if we created any routes
	if len(routes) == 0 {
		return nil, errors.New("could not create any routes with available trucks and delivery points")
	}

	log.Printf("Successfully created a total of %d routes across all categories", len(routes))
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

	// Create a map to track remaining pallets for each delivery point
	remainingPalletsMap := make(map[int]int)
	for _, point := range points {
		remainingPalletsMap[point.ID] = point.Pallets
	}

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

		// Create a new route with PointPallets map
		route := &models.Route{
			ID:             routeID,
			TruckID:        truck.ID,
			DeliveryDate:   deliveryDate,
			StartPoint:     startPoint.ID,
			DeliveryPoints: []int{startPoint.ID},
			PointPallets:   make(map[int]int),
			Status:         models.RouteStatus.Planned,
		}

		// Calculate how many pallets to assign to this route for the start point
		palletsToAssign := remainingPalletsMap[startPoint.ID]
		if palletsToAssign > truck.Capacity {
			palletsToAssign = truck.Capacity
		}

		// Add the pallets for the start point
		route.PointPallets[startPoint.ID] = palletsToAssign
		remainingPalletsMap[startPoint.ID] -= palletsToAssign

		// If this point has no more pallets, remove it from remaining
		if remainingPalletsMap[startPoint.ID] <= 0 {
			remaining = remaining[1:]
		}

		// Mark truck as unavailable
		truck.Available = false

		// Add more points using nearest neighbor algorithm
		currentPoint := startPoint
		currentCapacity := truck.Capacity - palletsToAssign

		for len(remaining) > 0 && currentCapacity > 0 {
			// Find the nearest point
			nearestIndex := -1
			nearestDistance := -1

			for i, point := range remaining {
				// Skip points with no remaining pallets
				if remainingPalletsMap[point.ID] <= 0 {
					continue
				}

				if distance, exists := currentPoint.Distances[point.ID]; exists {
					if nearestDistance == -1 || distance < nearestDistance {
						nearestDistance = distance
						nearestIndex = i
					}
				}
			}

			// If no suitable point found, break the loop
			if nearestIndex == -1 {
				break
			}

			// Get the next point
			nextPoint := remaining[nearestIndex]
			pointID := nextPoint.ID

			// Calculate how many pallets to assign from this point
			palletsToAssign := remainingPalletsMap[pointID]
			if palletsToAssign > currentCapacity {
				palletsToAssign = currentCapacity
			}

			// Only add the point if we can assign at least one pallet
			if palletsToAssign > 0 {
				// Check if this point is already in the route
				pointExists := false
				for _, id := range route.DeliveryPoints {
					if id == pointID {
						pointExists = true
						break
					}
				}

				// If point is not already in route, add it
				if !pointExists {
					route.DeliveryPoints = append(route.DeliveryPoints, pointID)
					route.PointPallets[pointID] = palletsToAssign
				} else {
					// If point is already in route, update its pallet count
					route.PointPallets[pointID] += palletsToAssign
				}

				// Update remaining pallets and capacity
				remainingPalletsMap[pointID] -= palletsToAssign
				currentCapacity -= palletsToAssign

				// Update current point
				currentPoint = nextPoint

				// If this point has no more pallets, remove it from remaining
				if remainingPalletsMap[pointID] <= 0 {
					remaining = append(remaining[:nearestIndex], remaining[nearestIndex+1:]...)
				}

				// Break if we've reached a reasonable number of points per route
				if len(route.DeliveryPoints) >= 5 {
					break
				}
			} else {
				// If no pallets can be assigned, remove the point from consideration
				remaining = append(remaining[:nearestIndex], remaining[nearestIndex+1:]...)
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

// generateSplitRoutesForCategory creates routes for delivery points by splitting them across multiple smaller trucks
func (s *RouteOptimizerService) generateSplitRoutesForCategory(
	ctx context.Context,
	points []*models.DeliveryPoint,
	availableTrucks []*models.Truck,
	alternativeCapacities []int,
	pointMap map[int]*models.DeliveryPoint,
	deliveryDate time.Time,
) ([]*models.Route, error) {
	routes := make([]*models.Route, 0)

	// Find trucks with alternative capacities
	suitableTrucks := make([]*models.Truck, 0)
	for _, capacity := range alternativeCapacities {
		for _, truck := range availableTrucks {
			if truck.Capacity == capacity && truck.Available {
				// Check if the truck can access the delivery points
				canAccess := true
				for _, point := range points {
					if !point.CanAcceptTruck(truck.Capacity) {
						canAccess = false
						break
					}
				}
				if canAccess {
					suitableTrucks = append(suitableTrucks, truck)
				}
			}
		}
	}

	// If no suitable trucks found with alternative capacities, try using any available truck
	if len(suitableTrucks) == 0 {
		log.Println("No trucks with specified alternative capacities available, trying any available truck")
		// Try to find any available truck that can handle at least some of the delivery points
		for _, truck := range availableTrucks {
			if truck.Available {
				// Check if this truck can access at least one delivery point
				for _, point := range points {
					if point.CanAcceptTruck(truck.Capacity) {
						suitableTrucks = append(suitableTrucks, truck)
						break
					}
				}
			}
		}
	}

	// If still no suitable trucks, return a more descriptive error
	if len(suitableTrucks) == 0 {
		return nil, errors.New("no suitable trucks available for delivery points - please add more trucks or adjust delivery point requirements")
	}

	// Sort trucks by capacity (largest first) to prioritize larger trucks
	sort.Slice(suitableTrucks, func(i, j int) bool {
		return suitableTrucks[i].Capacity > suitableTrucks[j].Capacity
	})

	// Sort points by number of pallets (descending)
	sort.Slice(points, func(i, j int) bool {
		return points[i].Pallets > points[j].Pallets
	})

	// Create a map to track remaining pallets for each delivery point
	remainingPalletsMap := make(map[int]int)
	for _, point := range points {
		remainingPalletsMap[point.ID] = point.Pallets
	}

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

		// Filter remaining points to only those that can be accessed by this truck
		accessiblePoints := make([]*models.DeliveryPoint, 0)
		for _, point := range remaining {
			if point.CanAcceptTruck(truck.Capacity) && remainingPalletsMap[point.ID] > 0 {
				accessiblePoints = append(accessiblePoints, point)
			}
		}

		// If no points can be accessed by this truck, move to the next truck
		if len(accessiblePoints) == 0 {
			truckIndex++
			continue
		}

		// Start with the first accessible point
		startPoint := accessiblePoints[0]

		// Create a new route with PointPallets map
		route := &models.Route{
			ID:             routeID,
			TruckID:        truck.ID,
			DeliveryDate:   deliveryDate,
			StartPoint:     startPoint.ID,
			DeliveryPoints: []int{startPoint.ID},
			PointPallets:   make(map[int]int),
			Status:         models.RouteStatus.Planned,
		}

		// Calculate how many pallets to assign to this route for the start point
		palletsToAssign := remainingPalletsMap[startPoint.ID]
		if palletsToAssign > truck.Capacity {
			palletsToAssign = truck.Capacity
		}

		// Add the pallets for the start point
		route.PointPallets[startPoint.ID] = palletsToAssign
		remainingPalletsMap[startPoint.ID] -= palletsToAssign

		// If this point has no more pallets, remove it from remaining
		if remainingPalletsMap[startPoint.ID] <= 0 {
			for i, point := range remaining {
				if point.ID == startPoint.ID {
					remaining = append(remaining[:i], remaining[i+1:]...)
					break
				}
			}
		}

		// Mark truck as unavailable
		truck.Available = false

		// Add more points using nearest neighbor algorithm
		currentPoint := startPoint
		currentCapacity := truck.Capacity - palletsToAssign

		for len(remaining) > 0 && currentCapacity > 0 {
			// Find the nearest point that can be accessed by this truck
			nearestIndex := -1
			nearestDistance := -1

			for i, point := range remaining {
				// Skip points that can't be accessed by this truck or have no remaining pallets
				if !point.CanAcceptTruck(truck.Capacity) || remainingPalletsMap[point.ID] <= 0 {
					continue
				}

				if distance, exists := currentPoint.Distances[point.ID]; exists {
					if nearestDistance == -1 || distance < nearestDistance {
						nearestDistance = distance
						nearestIndex = i
					}
				}
			}

			// If no accessible point found, break the loop
			if nearestIndex == -1 {
				break
			}

			// Get the next point
			nextPoint := remaining[nearestIndex]
			pointID := nextPoint.ID

			// Calculate how many pallets to assign from this point
			palletsToAssign := remainingPalletsMap[pointID]
			if palletsToAssign > currentCapacity {
				palletsToAssign = currentCapacity
			}

			// Only add the point if we can assign at least one pallet
			if palletsToAssign > 0 {
				// Check if this point is already in the route
				pointExists := false
				for _, id := range route.DeliveryPoints {
					if id == pointID {
						pointExists = true
						break
					}
				}

				// If point is not already in route, add it
				if !pointExists {
					route.DeliveryPoints = append(route.DeliveryPoints, pointID)
					route.PointPallets[pointID] = palletsToAssign
				} else {
					// If point is already in route, update its pallet count
					route.PointPallets[pointID] += palletsToAssign
				}

				// Update remaining pallets and capacity
				remainingPalletsMap[pointID] -= palletsToAssign
				currentCapacity -= palletsToAssign

				// Update current point
				currentPoint = nextPoint

				// If this point has no more pallets, remove it from remaining
				if remainingPalletsMap[pointID] <= 0 {
					remaining = append(remaining[:nearestIndex], remaining[nearestIndex+1:]...)
				}

				// Break if we've reached a reasonable number of points per route
				if len(route.DeliveryPoints) >= 5 {
					break
				}
			} else {
				// If no pallets can be assigned, remove the point from consideration
				remaining = append(remaining[:nearestIndex], remaining[nearestIndex+1:]...)
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
	if len(routes) == 0 {
		return errors.New("no routes to save")
	}

	// Get all delivery points for distance and cost calculations
	deliveryPoints, err := s.DeliveryRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get delivery points: %w", err)
	}

	// Create a map of delivery points for easier access
	pointMap := make(map[int]*models.DeliveryPoint)
	for _, point := range deliveryPoints {
		pointMap[point.ID] = point
	}

	// Create a map to track remaining pallets for each delivery point
	remainingPalletsMap := make(map[int]int)
	for _, point := range deliveryPoints {
		if point.Pallets > 0 {
			remainingPalletsMap[point.ID] = point.Pallets
		}
	}

	// Sort routes by truck capacity (largest first) to prioritize larger trucks
	sort.Slice(routes, func(i, j int) bool {
		truckI, _ := s.TruckRepo.GetByID(ctx, routes[i].TruckID)
		truckJ, _ := s.TruckRepo.GetByID(ctx, routes[j].TruckID)

		// If we couldn't get the truck info, put that route last
		if truckI == nil {
			return false
		}
		if truckJ == nil {
			return true
		}

		return truckI.Capacity > truckJ.Capacity
	})

	// Save routes, avoiding duplicates and tracking pallets
	savedRoutes := make([]*models.Route, 0)
	processedPoints := make(map[int]bool)

	for _, route := range routes {
		// Skip routes that have already been processed
		skipRoute := true
		for _, pointID := range route.DeliveryPoints {
			if !processedPoints[pointID] && remainingPalletsMap[pointID] > 0 {
				skipRoute = false
				break
			}
		}

		if skipRoute {
			log.Printf("Skipping route as all its delivery points have already been processed")
			continue
		}

		// Check if this route still has valid delivery points with remaining pallets
		validDeliveryPoints := make([]int, 0)
		pointPallets := make(map[int]int)

		// Get the truck for this route
		truck, err := s.TruckRepo.GetByID(ctx, route.TruckID)
		if err != nil {
			log.Printf("Warning: Failed to get truck %d: %v", route.TruckID, err)
			continue
		}

		// Calculate total capacity for this truck
		remainingCapacity := truck.Capacity

		for _, pointID := range route.DeliveryPoints {
			// Skip points that have no remaining pallets or have been fully processed
			if remainingPalletsMap[pointID] <= 0 || processedPoints[pointID] {
				continue
			}

			// Calculate how many pallets to assign from this point
			palletsToAssign := remainingPalletsMap[pointID]
			if palletsToAssign > remainingCapacity {
				palletsToAssign = remainingCapacity
			}

			if palletsToAssign > 0 {
				validDeliveryPoints = append(validDeliveryPoints, pointID)
				pointPallets[pointID] = palletsToAssign
				remainingPalletsMap[pointID] -= palletsToAssign
				remainingCapacity -= palletsToAssign

				// If all pallets for this point have been assigned, mark it as processed
				if remainingPalletsMap[pointID] <= 0 {
					processedPoints[pointID] = true
				}
			}

			// If truck is full, stop adding points
			if remainingCapacity <= 0 {
				break
			}
		}

		// Skip routes with no valid delivery points
		if len(validDeliveryPoints) == 0 {
			log.Printf("Skipping route with no valid delivery points")
			continue
		}

		// Update the route with valid delivery points and pallet counts
		route.DeliveryPoints = validDeliveryPoints
		route.PointPallets = pointPallets

		// Recalculate distance and cost with updated delivery points
		route.CalculateTotalDistance(pointMap)
		route.CalculateTotalCost(truck)

		// Save the route
		if err := s.RouteRepo.Create(ctx, route); err != nil {
			return fmt.Errorf("failed to create route: %w", err)
		}

		savedRoutes = append(savedRoutes, route)

		// Update truck availability
		truck.Available = false
		if err := s.TruckRepo.Update(ctx, truck); err != nil {
			log.Printf("Warning: Failed to update truck %d availability: %v", truck.ID, err)
			// Continue processing other routes
		}
	}

	log.Printf("Saved %d routes out of %d generated routes", len(savedRoutes), len(routes))

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

// DeleteRoute deletes a route and updates the truck's availability status
func (s *RouteOptimizerService) DeleteRoute(ctx context.Context, routeID int) error {
	// Get the route to find the associated truck
	route, err := s.RouteRepo.GetByID(ctx, routeID)
	if err != nil {
		return err
	}

	// Store the truck ID before deleting the route
	truckID := route.TruckID

	// Delete the route
	if err := s.RouteRepo.Delete(ctx, routeID); err != nil {
		return err
	}

	// Update the truck's availability status
	truck, err := s.TruckRepo.GetByID(ctx, truckID)
	if err != nil {
		log.Printf("Warning: Failed to get truck %d after deleting route: %v", truckID, err)
		return nil // Don't fail the route deletion if truck update fails
	}

	// Set the truck as available again
	truck.Available = true
	if err := s.TruckRepo.Update(ctx, truck); err != nil {
		log.Printf("Warning: Failed to update truck %d availability: %v", truck.ID, err)
	}

	return nil
}

// DeleteAllRoutes deletes all routes and resets truck availability
func (s *RouteOptimizerService) DeleteAllRoutes(ctx context.Context) error {
	// Get all routes to find associated trucks
	routes, err := s.RouteRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get routes: %w", err)
	}

	// Create a map of truck IDs to avoid duplicates
	truckIDs := make(map[int]bool)
	for _, route := range routes {
		truckIDs[route.TruckID] = true
	}

	// Delete all routes
	if err := s.RouteRepo.DeleteAll(ctx); err != nil {
		return fmt.Errorf("failed to delete routes: %w", err)
	}

	// Reset availability for all trucks that were in use
	for truckID := range truckIDs {
		truck, err := s.TruckRepo.GetByID(ctx, truckID)
		if err != nil {
			log.Printf("Warning: Failed to get truck %d: %v", truckID, err)
			continue
		}

		truck.Available = true
		if err := s.TruckRepo.Update(ctx, truck); err != nil {
			log.Printf("Warning: Failed to update truck %d availability: %v", truck.ID, err)
		}
	}

	// As a fallback, reset all trucks to available state
	allTrucks, err := s.TruckRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Warning: Failed to get all trucks: %v", err)
		return nil
	}

	for _, truck := range allTrucks {
		if !truck.Available {
			truck.Available = true
			if err := s.TruckRepo.Update(ctx, truck); err != nil {
				log.Printf("Warning: Failed to reset truck %d availability: %v", truck.ID, err)
			}
		}
	}

	return nil
}
