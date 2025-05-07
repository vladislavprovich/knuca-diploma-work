package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
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

	// If no routes were generated, return an error
	if len(routes) == 0 {
		return nil, errors.New("failed to generate any routes")
	}

	// Set route IDs and other properties
	nextRouteID, err := s.RouteRepo.GetNextRouteID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get next route ID: %w", err)
	}

	for i, route := range routes {
		route.ID = nextRouteID + i
		route.DeliveryDate = deliveryDate
		route.Status = models.RouteStatus.Planned
	}

	return routes, nil
}

// generateRoutesForCategory creates routes for delivery points of a specific category
func (s *RouteOptimizerService) generateRoutesForCategory(
	ctx context.Context,
	points []*models.DeliveryPoint,
	availableTrucks []*models.Truck,
	preferredCapacity int,
	pointMap map[int]*models.DeliveryPoint,
	deliveryDate time.Time,
) ([]*models.Route, error) {
	// Filter trucks by preferred capacity
	var suitableTrucks []*models.Truck
	for _, truck := range availableTrucks {
		if truck.Capacity == preferredCapacity {
			suitableTrucks = append(suitableTrucks, truck)
		}
	}

	if len(suitableTrucks) == 0 {
		return nil, fmt.Errorf("no suitable trucks with capacity %d available", preferredCapacity)
	}

	// Sort points by number of pallets (descending)
	sort.Slice(points, func(i, j int) bool {
		return points[i].Pallets > points[j].Pallets
	})

	// Generate routes
	routes := make([]*models.Route, 0)
	remainingPoints := make([]*models.DeliveryPoint, len(points))
	copy(remainingPoints, points)

	truckIndex := 0
	for len(remainingPoints) > 0 && truckIndex < len(suitableTrucks) {
		truck := suitableTrucks[truckIndex]

		// Create a new route
		route := &models.Route{
			TruckID:        truck.ID,
			DeliveryDate:   deliveryDate,
			DeliveryPoints: []int{},
			PointPallets:   make(map[int]int),
			Status:         models.RouteStatus.Planned,
			StartPoint:     1000, // Set warehouse as default start point
		}

		// Calculate remaining capacity for this truck
		remainingCapacity := truck.Capacity

		// Assign points to the route using a greedy approach
		for len(remainingPoints) > 0 && remainingCapacity > 0 {
			// Find the nearest point that fits in the remaining capacity
			nearestIndex := -1
			nearestDistance := -1

			for i, point := range remainingPoints {
				// Skip points that don't fit in the remaining capacity
				if point.Pallets > remainingCapacity {
					continue
				}

				// If this is the first point, select it
				if len(route.DeliveryPoints) == 0 {
					nearestIndex = i
					break
				}

				// Calculate distance from the last point in the route
				lastPointID := route.DeliveryPoints[len(route.DeliveryPoints)-1]
				lastPoint := pointMap[lastPointID]

				distance := lastPoint.Distances[point.ID]
				if nearestIndex == -1 || distance < nearestDistance {
					nearestIndex = i
					nearestDistance = distance
				}
			}

			// If no suitable point was found, break the loop
			if nearestIndex == -1 {
				break
			}

			// Add the point to the route
			point := remainingPoints[nearestIndex]

			// Check if we can assign all pallets from this point
			if point.Pallets <= remainingCapacity {
				// Add all pallets from this point
				route.DeliveryPoints = append(route.DeliveryPoints, point.ID)
				route.PointPallets[point.ID] = point.Pallets
				remainingCapacity -= point.Pallets

				// Remove the point from consideration
				remaining := make([]*models.DeliveryPoint, 0, len(remainingPoints)-1)
				for i, p := range remainingPoints {
					if i != nearestIndex {
						remaining = append(remaining, p)
					}
				}
				remainingPoints = remaining
			} else {
				// If no pallets can be assigned, remove the point from consideration
				remaining := make([]*models.DeliveryPoint, 0, len(remainingPoints)-1)
				for i, p := range remainingPoints {
					if i != nearestIndex {
						remaining = append(remaining, p)
					}
				}
				remainingPoints = remaining
			}
		}

		// Calculate total distance using real-world distances and cost
		// First calculate using the OSRM API for real-world distances
		if err := s.calculateRealRouteDistance(route, pointMap); err != nil {
			// Fallback to the simpler distance calculation if OSRM fails
			route.CalculateTotalDistance(pointMap)
		}
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

	// Save all routes without skipping any
	savedRoutes := make([]*models.Route, 0)

	for _, route := range routes {
		// Get the truck for this route
		truck, err := s.TruckRepo.GetByID(ctx, route.TruckID)
		if err != nil {
			log.Printf("Warning: Failed to get truck %d: %v", route.TruckID, err)
			continue
		}

		// Calculate accurate road distances using OSRM API
		if err := s.calculateRealRouteDistance(route, pointMap); err != nil {
			log.Printf("Warning: Failed to calculate real route distance: %v. Using default distance calculation.", err)
			// Fallback to the default distance calculation
			route.CalculateTotalDistance(pointMap)
			log.Printf("Route %d: Fallback distance calculation: %.1f km with %d delivery points",
				route.ID, route.TotalDistance, len(route.DeliveryPoints))
		}

		// Recalculate cost based on the updated distance
		route.CalculateTotalCost(truck)

		// Save the route as is without filtering delivery points
		if err := s.RouteRepo.Create(ctx, route); err != nil {
			log.Printf("Error creating route: %v", err)
			return fmt.Errorf("failed to create route: %w", err)
		}

		// Add to saved routes list
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

	// Calculate accurate road distances using OSRM API
	if err := s.calculateRealRouteDistance(route, pointMap); err != nil {
		log.Printf("Warning: Failed to calculate real route distance: %v. Using default distance calculation.", err)
		// Fallback to the default distance calculation
		route.CalculateTotalDistance(pointMap)
		log.Printf("Route %d: Fallback distance calculation: %.1f km with %d delivery points",
			route.ID, route.TotalDistance, len(route.DeliveryPoints))
	}

	// Recalculate cost based on the updated distance
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

// calculateRealRouteDistance calculates the actual road distance for a route using the OSRM API
func (s *RouteOptimizerService) calculateRealRouteDistance(route *models.Route, pointMap map[int]*models.DeliveryPoint) error {
	if len(route.DeliveryPoints) == 0 {
		route.TotalDistance = 0
		return nil
	}

	// Always use warehouse as start point to match map calculation
	warehouseID := 1000 // Vyshneve warehouse ID

	// Create an array of waypoints for the route
	waypoints := make([][]float64, 0)

	// Add warehouse as start point to match the map calculation
	if warehouse, exists := pointMap[warehouseID]; exists {
		waypoints = append(waypoints, []float64{warehouse.Longitude, warehouse.Latitude})
	}

	// Add all delivery points in order
	for _, pointID := range route.DeliveryPoints {
		if point, exists := pointMap[pointID]; exists {
			waypoints = append(waypoints, []float64{point.Longitude, point.Latitude})
		}
	}

	// Always add warehouse as ending point to complete the route
	if warehouse, exists := pointMap[warehouseID]; exists {
		waypoints = append(waypoints, []float64{warehouse.Longitude, warehouse.Latitude})
	}

	// Calculate total real distance using OSRM API
	totalRealDistance := 0.0

	// Process route segments (OSRM has a limit on number of waypoints)
	for i := 0; i < len(waypoints)-1; i++ {
		start := waypoints[i]
		end := waypoints[i+1]

		// Use OSRM public API to get route between two points
		url := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false",
			start[0], start[1], end[0], end[1])

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Parse the JSON response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}

		// Extract the distance from the response
		if routes, ok := result["routes"].([]interface{}); ok && len(routes) > 0 {
			if routeData, ok := routes[0].(map[string]interface{}); ok {
				if distance, ok := routeData["distance"].(float64); ok {
					// Convert meters to kilometers and round to one decimal place
					distanceKm := math.Round((distance/1000)*10) / 10
					totalRealDistance += distanceKm
				}
			}
		}

		// Add a small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	// Store the calculated distance in the route object with one decimal place precision
	route.TotalDistance = math.Round(totalRealDistance*10) / 10

	// Log the calculated distance for debugging
	log.Printf("Route %d: Calculated real distance: %.1f km with %d delivery points",
		route.ID, route.TotalDistance, len(route.DeliveryPoints))

	return nil
}

// DeleteAllRoutes deletes all routes and resets truck availability
func (s *RouteOptimizerService) DeleteAllRoutes(ctx context.Context) error {
	// Get all routes to find associated trucks
	routes, err := s.RouteRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Create a set of truck IDs to avoid duplicates
	truckIDs := make(map[int]bool)
	for _, route := range routes {
		truckIDs[route.TruckID] = true
	}

	// Delete all routes
	if err := s.RouteRepo.DeleteAll(ctx); err != nil {
		return err
	}

	// Reset availability for all trucks that were used in routes
	for truckID := range truckIDs {
		truck, err := s.TruckRepo.GetByID(ctx, truckID)
		if err != nil {
			log.Printf("Warning: Failed to get truck %d after deleting routes: %v", truckID, err)
			continue
		}

		// Set the truck as available again
		truck.Available = true
		if err := s.TruckRepo.Update(ctx, truck); err != nil {
			log.Printf("Warning: Failed to update truck %d availability: %v", truck.ID, err)
		}
	}

	return nil
}

// generateSplitRoutesForCategory creates routes for delivery points using alternative truck capacities
func (s *RouteOptimizerService) generateSplitRoutesForCategory(
	ctx context.Context,
	points []*models.DeliveryPoint,
	availableTrucks []*models.Truck,
	alternativeCapacities []int,
	pointMap map[int]*models.DeliveryPoint,
	deliveryDate time.Time,
) ([]*models.Route, error) {
	// Try each alternative capacity
	var allRoutes []*models.Route

	for _, capacity := range alternativeCapacities {
		// Filter trucks by capacity
		var suitableTrucks []*models.Truck
		for _, truck := range availableTrucks {
			if truck.Capacity == capacity {
				suitableTrucks = append(suitableTrucks, truck)
			}
		}

		if len(suitableTrucks) == 0 {
			log.Printf("No trucks with capacity %d available, trying next alternative", capacity)
			continue
		}

		// Generate routes with these trucks
		routes, err := s.generateRoutesForCategory(ctx, points, suitableTrucks, capacity, pointMap, deliveryDate)
		if err != nil {
			log.Printf("Failed to generate routes with %d-pallet trucks: %v", capacity, err)
			continue
		}

		// Add these routes to the result
		allRoutes = append(allRoutes, routes...)

		// Update the points list to remove points that have been fully assigned
		var remainingPoints []*models.DeliveryPoint
		for _, point := range points {
			// Check if this point is in any of the routes
			assigned := false
			for _, route := range routes {
				for _, pointID := range route.DeliveryPoints {
					if pointID == point.ID {
						assigned = true
						break
					}
				}
				if assigned {
					break
				}
			}

			if !assigned {
				remainingPoints = append(remainingPoints, point)
			}
		}

		// If no points remain, we're done
		if len(remainingPoints) == 0 {
			break
		}

		// Update points for the next iteration
		points = remainingPoints
	}

	if len(allRoutes) == 0 {
		return nil, errors.New("failed to generate routes with alternative truck capacities")
	}

	return allRoutes, nil
}
