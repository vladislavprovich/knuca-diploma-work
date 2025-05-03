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
	_, exists := CategoryPalletCapacity[dp.Category]
	if !exists {
		return false
	}

	// Implement category-based routing rules with fallback for flexibility:
	// 1. Blue category (33) can be accessed by 33-pallet trucks primarily
	//    but can accept 18-pallet trucks if necessary (may require multiple trips)
	// 2. Green category (18) can be accessed by 18-pallet trucks and smaller
	// 3. Yellow category (15) can be accessed by 15-pallet trucks and smaller
	// 4. Purple category (10) can be accessed by 10-pallet trucks
	switch dp.Category {
	case "Blue":
		// Blue points prefer 33-pallet trucks but can use 18-pallet trucks if needed
		return truckCapacity == 33 || truckCapacity == 18
	case "Green":
		return truckCapacity <= 18
	case "Yellow":
		return truckCapacity <= 15
	case "Purple":
		return truckCapacity <= 10
	default:
		return false
	}
}
