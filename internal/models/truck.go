package models

// Truck represents a delivery vehicle in the fleet
type Truck struct {
	ID           int    `json:"id" bson:"_id"`
	LicensePlate string `json:"license_plate" bson:"license_plate"`
	Brand        string `json:"brand" bson:"brand"`
	Capacity     int    `json:"capacity" bson:"capacity"` // Number of pallets (10, 15, 18, or 33)
	DriverName   string `json:"driver_name" bson:"driver_name"`
	DriverSurname string `json:"driver_surname" bson:"driver_surname"`
	HasTrailer   bool   `json:"has_trailer" bson:"has_trailer"` // Only true for 33-pallet trucks
	TrailerInfo  string `json:"trailer_info,omitempty" bson:"trailer_info,omitempty"`
	Available    bool   `json:"available" bson:"available"` // Whether the truck is available for assignment
}

// ValidTruckCapacities defines the allowed truck capacities in the system
var ValidTruckCapacities = []int{10, 15, 18, 33}

// TruckCostPerKm maps truck capacities to their cost per kilometer
var TruckCostPerKm = map[int]int{
	33: 62, // 33 pallets – 62 units/km
	18: 58, // 18 pallets – 58 units/km
	15: 56, // 15 pallets – 56 units/km
	10: 52, // 10 pallets – 52 units/km
}

// IsValidCapacity checks if the provided capacity is valid for a truck
func IsValidCapacity(capacity int) bool {
	for _, validCapacity := range ValidTruckCapacities {
		if capacity == validCapacity {
			return true
		}
	}
	return false
}

// GetCostPerKm returns the cost per kilometer for a truck with the given capacity
func GetCostPerKm(capacity int) (int, bool) {
	cost, exists := TruckCostPerKm[capacity]
	return cost, exists
}

// ValidateTrailer ensures that only 33-pallet trucks have trailers
func (t *Truck) ValidateTrailer() bool {
	if t.Capacity == 33 {
		return t.HasTrailer
	}
	return !t.HasTrailer // Other trucks should not have trailers
}