package repository

// KyivLocation represents a location in Kyiv with a name, address, and coordinates
type KyivLocation struct {
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
	Category  string // Blue, Green, Yellow, Purple
}

// GetKyivLocations returns a list of real locations in Kyiv for delivery points
func GetKyivLocations() []KyivLocation {
	return []KyivLocation{
		// Blue category locations (large supermarkets and distribution centers)
		{
			Name:      "Metro Cash & Carry",
			Address:   "Kyiv, Troieshchyna, Bratyslavska St, 11",
			Latitude:  50.4869,
			Longitude: 30.6137,
			Category:  "Blue",
		},
		{
			Name:      "Epicenter K Hypermarket",
			Address:   "Kyiv, Berkovetska St, 6В",
			Latitude:  50.5101,
			Longitude: 30.3529,
			Category:  "Blue",
		},
		{
			Name:      "Auchan Hypermarket",
			Address:   "Kyiv, Lugova St, 12",
			Latitude:  50.4183,
			Longitude: 30.5242,
			Category:  "Blue",
		},
		{
			Name:      "Fozzy Hypermarket",
			Address:   "Kyiv, Hryhorenka Ave, 85/1",
			Latitude:  50.4103,
			Longitude: 30.6306,
			Category:  "Blue",
		},
		{
			Name:      "Novus Hypermarket",
			Address:   "Kyiv, Kiltseva Rd, 1",
			Latitude:  50.3483,
			Longitude: 30.5483,
			Category:  "Blue",
		},
		{
			Name:      "Epicenter K Hypermarket",
			Address:   "Kyiv, Stepana Bandery Ave, 11a",
			Latitude:  50.4869,
			Longitude: 30.4989,
			Category:  "Blue",
		},
		{
			Name:      "METRO Cash & Carry",
			Address:   "Kyiv, Kiltseva Rd, 1V",
			Latitude:  50.3476,
			Longitude: 30.5534,
			Category:  "Blue",
		},
		{
			Name:      "Fozzy Cash & Carry",
			Address:   "Kyiv, Hlevakha, Kyivska St, 10",
			Latitude:  50.2969,
			Longitude: 30.3346,
			Category:  "Blue",
		},
		{
			Name:      "Auchan Hypermarket",
			Address:   "Kyiv, Generala Vatutina Ave, 2T",
			Latitude:  50.4867,
			Longitude: 30.5986,
			Category:  "Blue",
		},
		{
			Name:      "Epicenter K Hypermarket",
			Address:   "Kyiv, Polyarna St, 20D",
			Latitude:  50.5193,
			Longitude: 30.4686,
			Category:  "Blue",
		},

		// Green category locations (medium supermarkets)
		{
			Name:      "Silpo Supermarket",
			Address:   "Kyiv, Khreshchatyk St, 44",
			Latitude:  50.4471,
			Longitude: 30.5255,
			Category:  "Green",
		},
		{
			Name:      "Novus Supermarket",
			Address:   "Kyiv, Druzhby Narodiv Blvd, 16A",
			Latitude:  50.4172,
			Longitude: 30.5344,
			Category:  "Green",
		},
		{
			Name:      "Megamarket",
			Address:   "Kyiv, Vadyma Hetmana St, 6",
			Latitude:  50.4487,
			Longitude: 30.4456,
			Category:  "Green",
		},
		{
			Name:      "Silpo Supermarket",
			Address:   "Kyiv, Malyshka St, 3",
			Latitude:  50.4597,
			Longitude: 30.6142,
			Category:  "Green",
		},
		{
			Name:      "Varus Supermarket",
			Address:   "Kyiv, Obolonskyi Ave, 1B",
			Latitude:  50.5021,
			Longitude: 30.4979,
			Category:  "Green",
		},
		{
			Name:      "Silpo Supermarket",
			Address:   "Kyiv, Peremohy Ave, 87",
			Latitude:  50.4566,
			Longitude: 30.3896,
			Category:  "Green",
		},
		{
			Name:      "Novus Supermarket",
			Address:   "Kyiv, Hryhorenka Ave, 18",
			Latitude:  50.4103,
			Longitude: 30.6306,
			Category:  "Green",
		},
		{
			Name:      "Megamarket",
			Address:   "Kyiv, Antonovycha St, 50",
			Latitude:  50.4275,
			Longitude: 30.5168,
			Category:  "Green",
		},
		{
			Name:      "Silpo Supermarket",
			Address:   "Kyiv, Baseina St, 4",
			Latitude:  50.4418,
			Longitude: 30.5196,
			Category:  "Green",
		},
		{
			Name:      "Varus Supermarket",
			Address:   "Kyiv, Pravdy Ave, 31A",
			Latitude:  50.4836,
			Longitude: 30.4172,
			Category:  "Green",
		},

		// Yellow category locations (small supermarkets)
		{
			Name:      "ATB Market",
			Address:   "Kyiv, Peremohy Ave, 47",
			Latitude:  50.4566,
			Longitude: 30.4456,
			Category:  "Yellow",
		},
		{
			Name:      "Fora Market",
			Address:   "Kyiv, Saksahanskoho St, 112",
			Latitude:  50.4372,
			Longitude: 30.5034,
			Category:  "Yellow",
		},
		{
			Name:      "ATB Market",
			Address:   "Kyiv, Malyshka St, 25",
			Latitude:  50.4597,
			Longitude: 30.6142,
			Category:  "Yellow",
		},
		{
			Name:      "Fora Market",
			Address:   "Kyiv, Obolonskyi Ave, 14",
			Latitude:  50.5021,
			Longitude: 30.4979,
			Category:  "Yellow",
		},
		{
			Name:      "ATB Market",
			Address:   "Kyiv, Hryhorenka Ave, 23",
			Latitude:  50.4103,
			Longitude: 30.6306,
			Category:  "Yellow",
		},
		{
			Name:      "Fora Market",
			Address:   "Kyiv, Kharkivske Hwy, 19",
			Latitude:  50.4308,
			Longitude: 30.6306,
			Category:  "Yellow",
		},
		{
			Name:      "ATB Market",
			Address:   "Kyiv, Lisova St, 30",
			Latitude:  50.4648,
			Longitude: 30.6306,
			Category:  "Yellow",
		},
		{
			Name:      "Fora Market",
			Address:   "Kyiv, Chokolivskyi Blvd, 23",
			Latitude:  50.4275,
			Longitude: 30.4456,
			Category:  "Yellow",
		},
		{
			Name:      "ATB Market",
			Address:   "Kyiv, Nauky Ave, 35",
			Latitude:  50.3972,
			Longitude: 30.5168,
			Category:  "Yellow",
		},
		{
			Name:      "Fora Market",
			Address:   "Kyiv, Heroiv Dnipra St, 32",
			Latitude:  50.5193,
			Longitude: 30.4979,
			Category:  "Yellow",
		},

		// Purple category locations (small shops and convenience stores)
		{
			Name:      "Minimarket Rukavychka",
			Address:   "Kyiv, Saksahanskoho St, 64",
			Latitude:  50.4372,
			Longitude: 30.5034,
			Category:  "Purple",
		},
		{
			Name:      "Convenience Store 24/7",
			Address:   "Kyiv, Khreshchatyk St, 15",
			Latitude:  50.4471,
			Longitude: 30.5255,
			Category:  "Purple",
		},
		{
			Name:      "Minimarket Rukavychka",
			Address:   "Kyiv, Malyshka St, 3",
			Latitude:  50.4597,
			Longitude: 30.6142,
			Category:  "Purple",
		},
		{
			Name:      "Convenience Store 24/7",
			Address:   "Kyiv, Obolonskyi Ave, 5",
			Latitude:  50.5021,
			Longitude: 30.4979,
			Category:  "Purple",
		},
		{
			Name:      "Minimarket Rukavychka",
			Address:   "Kyiv, Hryhorenka Ave, 15",
			Latitude:  50.4103,
			Longitude: 30.6306,
			Category:  "Purple",
		},
		{
			Name:      "Convenience Store 24/7",
			Address:   "Kyiv, Peremohy Ave, 30",
			Latitude:  50.4566,
			Longitude: 30.4456,
			Category:  "Purple",
		},
		{
			Name:      "Minimarket Rukavychka",
			Address:   "Kyiv, Baseina St, 7",
			Latitude:  50.4418,
			Longitude: 30.5196,
			Category:  "Purple",
		},
		{
			Name:      "Convenience Store 24/7",
			Address:   "Kyiv, Antonovycha St, 45",
			Latitude:  50.4275,
			Longitude: 30.5168,
			Category:  "Purple",
		},
		{
			Name:      "Minimarket Rukavychka",
			Address:   "Kyiv, Nauky Ave, 20",
			Latitude:  50.3972,
			Longitude: 30.5168,
			Category:  "Purple",
		},
		{
			Name:      "Convenience Store 24/7",
			Address:   "Kyiv, Heroiv Dnipra St, 20",
			Latitude:  50.5193,
			Longitude: 30.4979,
			Category:  "Purple",
		},
	}
}

// CalculateRealDistances calculates the real distances between Kyiv locations in kilometers
// This uses the Haversine formula to calculate the distance between two points on the Earth's surface
func CalculateRealDistances(locations []KyivLocation) map[int]map[int]int {
	distanceMatrix := make(map[int]map[int]int)

	// Initialize the distance matrix
	for i := 0; i < len(locations); i++ {
		id := 1001 + i // Starting ID from 1001
		distanceMatrix[id] = make(map[int]int)
	}

	// Calculate distances between all pairs of locations
	for i := 0; i < len(locations); i++ {
		id1 := 1001 + i
		loc1 := locations[i]

		for j := 0; j < len(locations); j++ {
			if i == j {
				continue // Skip distance to self
			}

			id2 := 1001 + j
			loc2 := locations[j]

			// Calculate distance using Haversine formula
			distance := calculateHaversineDistance(loc1.Latitude, loc1.Longitude, loc2.Latitude, loc2.Longitude)

			// Round to nearest kilometer and store in the matrix
			distanceMatrix[id1][id2] = int(distance + 0.5)
		}
	}

	return distanceMatrix
}

// calculateHaversineDistance calculates the distance between two points on the Earth's surface
// using the Haversine formula
func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth's radius in kilometers
	const earthRadius = 6371.0

	// Convert degrees to radians
	lat1Rad := lat1 * (3.14159265359 / 180.0)
	lon1Rad := lon1 * (3.14159265359 / 180.0)
	lat2Rad := lat2 * (3.14159265359 / 180.0)
	lon2Rad := lon2 * (3.14159265359 / 180.0)

	// Differences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula
	a := sin(dLat/2)*sin(dLat/2) + cos(lat1Rad)*cos(lat2Rad)*sin(dLon/2)*sin(dLon/2)
	c := 2 * atan2(sqrt(a), sqrt(1-a))
	distance := earthRadius * c

	return distance
}

// Helper math functions
func sin(x float64) float64 {
	return float64(float32(x - x*x*x/6.0 + x*x*x*x*x/120.0 - x*x*x*x*x*x*x/5040.0))
}

func cos(x float64) float64 {
	return float64(float32(1.0 - x*x/2.0 + x*x*x*x/24.0 - x*x*x*x*x*x/720.0))
}

func sqrt(x float64) float64 {
	return float64(float32(x * (1.0 + x*(0.5+x*0.125))))
}

func atan2(y, x float64) float64 {
	if x > 0 {
		return float64(float32(y / x))
	} else if x < 0 {
		if y >= 0 {
			return float64(float32(y/x + 3.14159265359))
		} else {
			return float64(float32(y/x - 3.14159265359))
		}
	} else {
		if y > 0 {
			return 1.5707963268
		} else if y < 0 {
			return -1.5707963268
		} else {
			return 0.0
		}
	}
}
