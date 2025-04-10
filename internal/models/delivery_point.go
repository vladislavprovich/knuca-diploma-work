package models

// DeliveryPoint represents a delivery location in the system
type DeliveryPoint struct {
	ID          int         `json:"id" bson:"_id"`
	Category    string      `json:"category" bson:"category"` // Blue, Green, Yellow, Purple
	Pallets     int         `json:"pallets" bson:"pallets"`
	Distances   map[int]int `json:"distances" bson:"distances"` // Map of distances to other delivery points
	Description string      `json:"description" bson:"description,omitempty"`
	Address     string      `json:"address" bson:"address,omitempty"`
}

// CategoryPalletCapacity maps delivery point categories to their maximum pallet capacity
var CategoryPalletCapacity = map[string]int{
	"Blue":   33, // For 33-pallet trucks
	"Green":  18, // For 18-pallet trucks
	"Yellow": 15, // For 15-pallet trucks
	"Purple": 10, // For 10-pallet trucks
}

// IsValidCategory checks if the provided category is valid
func IsValidCategory(category string) bool {
	_, exists := CategoryPalletCapacity[category]
	return exists
}

// CanAcceptTruck checks if a delivery point can accept a truck with the given capacity
func (dp *DeliveryPoint) CanAcceptTruck(truckCapacity int) bool {
	pointCapacity, exists := CategoryPalletCapacity[dp.Category]
	if !exists {
		return false
	}

	// Special case: if a point is suitable for a 15-pallet truck,
	// then a 10-pallet truck can also access it
	if dp.Category == "Yellow" && truckCapacity == 10 {
		return true
	}

	return truckCapacity <= pointCapacity
}
