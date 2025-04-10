package models

import (
	"time"
)

// Route represents a delivery route with assigned truck and delivery points
type Route struct {
	ID             int       `json:"id" bson:"_id"`
	TruckID        int       `json:"truck_id" bson:"truck_id"`
	DeliveryDate   time.Time `json:"delivery_date" bson:"delivery_date"`
	StartPoint     int       `json:"start_point" bson:"start_point"`         // ID of the starting delivery point
	DeliveryPoints []int     `json:"delivery_points" bson:"delivery_points"` // IDs of delivery points in order
	TotalDistance  int       `json:"total_distance" bson:"total_distance"`   // Total route distance in km
	TotalCost      int       `json:"total_cost" bson:"total_cost"`           // Total cost of the route
	Status         string    `json:"status" bson:"status"`                   // Planned, In Progress, Completed, Cancelled
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

// RouteStatus defines the possible statuses for a route
var RouteStatus = struct {
	Planned    string
	InProgress string
	Completed  string
	Cancelled  string
}{
	Planned:    "Planned",
	InProgress: "In Progress",
	Completed:  "Completed",
	Cancelled:  "Cancelled",
}

// CalculateTotalDistance computes the total distance of the route based on the delivery points
func (r *Route) CalculateTotalDistance(deliveryPoints map[int]*DeliveryPoint) int {
	if len(r.DeliveryPoints) == 0 {
		return 0
	}

	totalDistance := 0
	currentPoint := r.StartPoint

	// Calculate distance from start point to first delivery point
	if firstPoint, exists := deliveryPoints[r.DeliveryPoints[0]]; exists {
		if startPoint, exists := deliveryPoints[currentPoint]; exists {
			if distance, exists := startPoint.Distances[firstPoint.ID]; exists {
				totalDistance += distance
			}
		}
	}

	// Calculate distances between consecutive delivery points
	for i := 0; i < len(r.DeliveryPoints)-1; i++ {
		currentPointID := r.DeliveryPoints[i]
		nextPointID := r.DeliveryPoints[i+1]

		if currentPoint, exists := deliveryPoints[currentPointID]; exists {
			if distance, exists := currentPoint.Distances[nextPointID]; exists {
				totalDistance += distance
			}
		}
	}

	r.TotalDistance = totalDistance
	return totalDistance
}

// CalculateTotalCost computes the total cost of the route based on distance and truck capacity
func (r *Route) CalculateTotalCost(truck *Truck) int {
	costPerKm, exists := GetCostPerKm(truck.Capacity)
	if !exists {
		return 0
	}

	r.TotalCost = r.TotalDistance * costPerKm
	return r.TotalCost
}

// ValidateRoute checks if the route is valid based on truck capacity and delivery point compatibility
func (r *Route) ValidateRoute(truck *Truck, deliveryPoints map[int]*DeliveryPoint) bool {
	// Check if truck exists and is available
	if truck == nil || !truck.Available {
		return false
	}

	// Check if all delivery points can accept the assigned truck
	for _, pointID := range r.DeliveryPoints {
		point, exists := deliveryPoints[pointID]
		if !exists {
			return false
		}

		if !point.CanAcceptTruck(truck.Capacity) {
			return false
		}
	}

	return true
}
