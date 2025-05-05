import React, { useState, useEffect } from 'react';
import { Box, Button, FormControl, InputLabel, Select, MenuItem, Typography, CircularProgress, Paper, Grid } from '@mui/material';
import { MapContainer, TileLayer, Marker, Popup, useMap, Polyline } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import './Map.css';
import L from 'leaflet';
import axios from 'axios';

// MapController component to handle map view updates
const MapController = ({ center, zoom }) => {
  const map = useMap();
  
  useEffect(() => {
    if (map) {
      map.setView(center, zoom);
    }
  }, [center, zoom, map]);
  
  return null;
};

// Fix for Leaflet marker icons
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
});

// Custom marker icons for different categories
const categoryIcons = {
  Warehouse: new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [30, 45], // Slightly larger for warehouse
    iconAnchor: [15, 45],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  }),
  Blue: new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-blue.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  }),
  Green: new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  }),
  Yellow: new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-yellow.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  }),
  Purple: new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-violet.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  }),
};

// Static data for when API is not available
const staticWarehouse = {
  id: 1000,
  category: 'Warehouse',
  pallets: 0,
  description: 'Central Warehouse',
  address: 'Vyshneve, Kyiv Region, Ukraine',
  latitude: 50.3833,
  longitude: 30.3667
};

// Static delivery points data based on kyiv_locations.go
const staticDeliveryPoints = [
  staticWarehouse,
  {
    id: 1001,
    category: 'Blue',
    pallets: 20,
    description: 'Metro Cash & Carry',
    address: 'Kyiv, Troieshchyna, Bratyslavska St, 11',
    latitude: 50.4869,
    longitude: 30.6137
  },
  {
    id: 1002,
    category: 'Blue',
    pallets: 25,
    description: 'Epicenter K Hypermarket',
    address: 'Kyiv, Berkovetska St, 6В',
    latitude: 50.5101,
    longitude: 30.3529
  },
  {
    id: 1003,
    category: 'Blue',
    pallets: 18,
    description: 'METRO Cash & Carry Odesa Highway',
    address: 'Kyiv, Kiltseva Rd, 1B',
    latitude: 50.3483,
    longitude: 30.4672
  },
  {
    id: 1011,
    category: 'Green',
    pallets: 15,
    description: 'Silpo Supermarket',
    address: 'Kyiv, Khreshchatyk St, 44',
    latitude: 50.4471,
    longitude: 30.5255
  },
  {
    id: 1012,
    category: 'Green',
    pallets: 12,
    description: 'Novus Supermarket',
    address: 'Kyiv, Druzhby Narodiv Blvd, 16A',
    latitude: 50.4172,
    longitude: 30.5344
  },
  {
    id: 1013,
    category: 'Green',
    pallets: 14,
    description: 'Silpo Supermarket Obolon',
    address: 'Kyiv, Obolonsky Ave, 19',
    latitude: 50.5021,
    longitude: 30.4984
  },
  {
    id: 1014,
    category: 'Green',
    pallets: 11,
    description: 'Silpo Supermarket Pozniaky',
    address: 'Kyiv, Mykhaila Hryshka St, 3A',
    latitude: 50.4008,
    longitude: 30.6306
  },
  {
    id: 1015,
    category: 'Green',
    pallets: 13,
    description: 'Megamarket Supermarket',
    address: 'Kyiv, Vadyma Hetmana St, 6',
    latitude: 50.4487,
    longitude: 30.4456
  },
  {
    id: 1021,
    category: 'Yellow',
    pallets: 8,
    description: 'ATB Market',
    address: 'Kyiv, Peremohy Ave, 47',
    latitude: 50.4566,
    longitude: 30.4456
  },
  {
    id: 1022,
    category: 'Yellow',
    pallets: 10,
    description: 'Fora Market',
    address: 'Kyiv, Saksahanskoho St, 112',
    latitude: 50.4372,
    longitude: 30.5034
  },
  {
    id: 1023,
    category: 'Yellow',
    pallets: 7,
    description: 'ATB Market Sviatoshyn',
    address: 'Kyiv, Peremohy Ave, 87',
    latitude: 50.4577,
    longitude: 30.3905
  },
  {
    id: 1024,
    category: 'Yellow',
    pallets: 9,
    description: 'ATB Market Darnytsia',
    address: 'Kyiv, Kharkivske Highway, 168',
    latitude: 50.4315,
    longitude: 30.6841
  },
  {
    id: 1025,
    category: 'Yellow',
    pallets: 6,
    description: 'Fora Market Pechersk',
    address: 'Kyiv, Lesi Ukrainky Blvd, 24',
    latitude: 50.4237,
    longitude: 30.5362
  },
  {
    id: 1031,
    category: 'Purple',
    pallets: 5,
    description: 'Minimarket Rukavychka',
    address: 'Kyiv, Saksahanskoho St, 64',
    latitude: 50.4372,
    longitude: 30.5034
  },
  {
    id: 1032,
    category: 'Purple',
    pallets: 4,
    description: 'Convenience Store 24/7',
    address: 'Kyiv, Khreshchatyk St, 15',
    latitude: 50.4471,
    longitude: 30.5255
  }
];

// Custom style for the map container - absolute positioning to fill the entire viewport
const mapContainerStyle = {
  position: 'absolute',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  height: '100vh',
  width: '100%',
  backgroundColor: '#ffffff',
  padding: 0,
  margin: 0,
  zIndex: 1000 // Ensure map is above other elements
};

// Custom style for the map wrapper - position relative to contain the absolute positioned map
const mapWrapperStyle = {
  position: 'relative',
  height: '100vh',
  width: '100%',
  padding: 0,
  margin: 0,
  overflow: 'hidden',
  backgroundColor: '#ffffff'
};

const Map = () => {
  const [mapCenter, setMapCenter] = useState([50.3833, 30.3667]); // Default center (Vyshneve warehouse)
  const [mapZoom, setMapZoom] = useState(10);
  const [deliveryPoints, setDeliveryPoints] = useState([]);
  const [routes, setRoutes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedRoute, setSelectedRoute] = useState(null);
  const [routePolylines, setRoutePolylines] = useState([]);
  const [routeSegments, setRouteSegments] = useState([]);
  
  // Fetch delivery points and routes from API
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        // Fetch delivery points
        const pointsResponse = await axios.get('/api/delivery-points');
        setDeliveryPoints(pointsResponse.data);
        
        // Fetch routes
        const routesResponse = await axios.get('/api/routes');
        setRoutes(routesResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
        // Fall back to static data if API fails
        setDeliveryPoints(staticDeliveryPoints);
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, []);
  
  // Generate polylines when a route is selected
  useEffect(() => {
    if (!selectedRoute) {
      setRoutePolylines([]);
      setRouteSegments([]);
      return;
    }
    
    // Find all delivery points in the route
    const pointsMap = {};
    const segmentsInfo = [];
    
    // Create a map of all delivery points by ID for quick lookup
    deliveryPoints.forEach(point => {
      pointsMap[point.id] = point;
    });
    
    // Add warehouse as starting point
    const warehouse = pointsMap[1000] || staticWarehouse;
    
    // Create an array of waypoints for the route
    const waypoints = [];
    waypoints.push([warehouse.longitude, warehouse.latitude]); // Start with warehouse
    
    // Add all delivery points in the route in order
    const orderedPoints = [];
    orderedPoints.push(warehouse);
    
    if (selectedRoute.delivery_points && selectedRoute.delivery_points.length > 0) {
      selectedRoute.delivery_points.forEach(pointId => {
        const point = pointsMap[pointId];
        if (point) {
          waypoints.push([point.longitude, point.latitude]);
          orderedPoints.push(point);
        }
      });
    }
    
    // Add warehouse as ending point to complete the route
    waypoints.push([warehouse.longitude, warehouse.latitude]);
    orderedPoints.push(warehouse);
    
    // Use OSRM API to get road-following routes with accurate distances
    const fetchRoadRoutes = async () => {
      try {
        // Create temporary array to hold all polyline segments
        const allSegments = [];
        const routeSegmentsData = [];
        let totalRealDistance = 0;
        
        // Process route segments (OSRM has a limit on number of waypoints)
        for (let i = 0; i < waypoints.length - 1; i++) {
          const start = waypoints[i];
          const end = waypoints[i + 1];
          const startPoint = orderedPoints[i];
          const endPoint = orderedPoints[i + 1];
          
          // Use OSRM public API to get route between two points
          const response = await fetch(
            `https://router.project-osrm.org/route/v1/driving/${start[0]},${start[1]};${end[0]},${end[1]}?overview=full&geometries=geojson`
          );
          
          const data = await response.json();
          
          if (data.routes && data.routes.length > 0) {
            // Get the real distance from OSRM in kilometers (OSRM returns meters)
            const realDistance = Math.round((data.routes[0].distance / 1000) * 10) / 10;
            totalRealDistance += realDistance;
            
            // OSRM returns coordinates as [longitude, latitude], but Leaflet needs [latitude, longitude]
            const coordinates = data.routes[0].geometry.coordinates.map(coord => [
              coord[1], coord[0]
            ]);
            
            allSegments.push(coordinates);
            
            // Get the midpoint of this segment for the distance label
            const midIndex = Math.floor(coordinates.length / 2);
            const midpoint = coordinates[midIndex];
            
            // Add segment info with coordinates and real distance
            routeSegmentsData.push({
              coordinates,
              midpoint,
              distance: realDistance,
              startName: startPoint.description || 'Warehouse',
              endName: endPoint.description || 'Warehouse',
              isReturnTrip: i === waypoints.length - 2 // Flag for return trip to warehouse
            });
          }
        }
        
        // Update the route's total distance with the real calculated distance
        if (selectedRoute) {
          selectedRoute.total_distance = Math.round(totalRealDistance * 10) / 10;
        }
        
        // Flatten all segments into a single array of points
        const flattenedRoute = allSegments.flat();
        setRoutePolylines(flattenedRoute);
        setRouteSegments(routeSegmentsData);
      } catch (error) {
        console.error('Error fetching road routes:', error);
        
        // Fallback to straight lines if OSRM fails
        const fallbackPoints = [];
        const fallbackSegments = [];
        let totalEstimatedDistance = 0;
        
        for (let i = 0; i < orderedPoints.length - 1; i++) {
          const startPoint = orderedPoints[i];
          const endPoint = orderedPoints[i + 1];
          
          const startCoord = [startPoint.latitude, startPoint.longitude];
          const endCoord = [endPoint.latitude, endPoint.longitude];
          const segmentCoords = [startCoord, endCoord];
          
          // Calculate straight-line distance using Haversine formula
          const R = 6371; // Earth's radius in km
          const dLat = (endPoint.latitude - startPoint.latitude) * Math.PI / 180;
          const dLon = (endPoint.longitude - startPoint.longitude) * Math.PI / 180;
          const a = 
            Math.sin(dLat/2) * Math.sin(dLat/2) +
            Math.cos(startPoint.latitude * Math.PI / 180) * Math.cos(endPoint.latitude * Math.PI / 180) * 
            Math.sin(dLon/2) * Math.sin(dLon/2);
          const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
          const distance = Math.round((R * c) * 10) / 10;
          
          // Add 30% to straight-line distance to approximate road distance
          const estimatedRoadDistance = Math.round(distance * 1.3 * 10) / 10;
          totalEstimatedDistance += estimatedRoadDistance;
          
          fallbackPoints.push(startCoord);
          if (i === orderedPoints.length - 2) {
            fallbackPoints.push(endCoord);
          }
          
          // Calculate midpoint for label
          const midpoint = [
            (startPoint.latitude + endPoint.latitude) / 2,
            (startPoint.longitude + endPoint.longitude) / 2
          ];
          
          fallbackSegments.push({
            coordinates: segmentCoords,
            midpoint,
            distance: estimatedRoadDistance,
            startName: startPoint.description || 'Warehouse',
            endName: endPoint.description || 'Warehouse',
            isReturnTrip: i === orderedPoints.length - 2 // Flag for return trip to warehouse
          });
        }
        
        // Update the route's total distance with the estimated distance
        if (selectedRoute) {
          selectedRoute.total_distance = Math.round(totalEstimatedDistance * 10) / 10;
        }
        
        setRoutePolylines(fallbackPoints);
        setRouteSegments(fallbackSegments);
      }
    };
    
    fetchRoadRoutes();
  }, [selectedRoute, deliveryPoints]);
  
  // Force map to load properly by using useEffect
  useEffect(() => {
    // This helps ensure the map container is properly sized when the component mounts
    const mapContainer = document.querySelector('.leaflet-container');
    if (mapContainer) {
      window.dispatchEvent(new Event('resize'));
    }
  }, []);
  
  // Handle route selection
  const handleRouteChange = (event) => {
    const routeId = event.target.value;
    if (routeId === '') {
      setSelectedRoute(null);
      return;
    }
    
    const route = routes.find(r => r.id.toString() === routeId);
    setSelectedRoute(route);
  };

  // Get points to display - use API data if available, otherwise use static data
  const pointsToDisplay = deliveryPoints.length > 0 ? deliveryPoints : staticDeliveryPoints;
  
  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }
  
  return (
    <Grid container sx={{ height: '100vh' }}>
      {/* Map container - takes 75% of the width */}
      <Grid item xs={12} md={9} sx={{ position: 'relative', height: '100vh', padding: 0 }}>
        <Box sx={mapWrapperStyle} className="map-wrapper">
          <div className="map-container-wrapper">
            <MapContainer 
              center={mapCenter} 
              zoom={mapZoom} 
              style={mapContainerStyle}
              zoomControl={true}
              attributionControl={true}
              className="map-container"
            >
              <MapController center={mapCenter} zoom={mapZoom} />
              <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                className="white-map-tiles"
              />
              
              {/* Display route polyline if a route is selected */}
              {routePolylines.length > 0 && (
                <>
                  {/* Render each route segment with different styling */}
                  {routeSegments.map((segment, index) => {
                    const isReturnTrip = segment.isReturnTrip;
                    return (
                      <Polyline 
                        key={`route-segment-${index}`}
                        positions={segment.coordinates}
                        color={isReturnTrip ? "#FF9800" : "#2196F3"} // Orange for return, Blue for outbound
                        weight={4}
                        opacity={0.8}
                        dashArray={isReturnTrip ? "10, 10" : null} // Dashed line for return trip
                        eventHandlers={{
                          click: (e) => {
                            const popup = L.popup()
                              .setLatLng(e.latlng)
                              .setContent(
                                `<div>
                                  <h4>${isReturnTrip ? "Return Trip" : "Delivery Segment"}</h4>
                                  <p><strong>From:</strong> ${segment.startName}</p>
                                  <p><strong>To:</strong> ${segment.endName}</p>
                                  <p><strong>Distance:</strong> ${segment.distance} km</p>
                                </div>`
                              )
                              .openOn(e.target._map);
                          }
                        }}
                      />
                    );
                  })}
                  
                  {/* Overall route information popup handler */}
                  <Polyline 
                    positions={routePolylines}
                    color="transparent"
                    weight={10}
                    opacity={0.01}
                    eventHandlers={{
                      click: (e) => {
                        if (selectedRoute) {
                          const popup = L.popup()
                            .setLatLng(e.latlng)
                            .setContent(
                              `<div>
                                <h4>Route Information</h4>
                                <p>Total Distance (including return): ${selectedRoute.total_distance} km</p>
                                <p>Delivery Points: ${selectedRoute.delivery_points ? selectedRoute.delivery_points.length : 0}</p>
                                <p>Total Pallets: ${selectedRoute.point_pallets ? Object.values(selectedRoute.point_pallets).reduce((sum, val) => sum + val, 0) : 0}</p>
                              </div>`
                            )
                            .openOn(e.target._map);
                        }
                      }
                    }}
                  />
                  {/* Render distance labels for each segment */}
                  {routeSegments.map((segment, index) => {
                    // Create a custom divIcon for the distance label
                    // Use a different style for the return trip
                    const isReturnTrip = segment.isReturnTrip;
                    const distanceIcon = L.divIcon({
                      className: 'distance-label',
                      html: `<div style="background-color: ${isReturnTrip ? '#ffecb3' : 'white'}; 
                                        padding: 3px 8px; 
                                        border-radius: 4px; 
                                        border: 1px solid ${isReturnTrip ? '#ff9800' : '#666'}; 
                                        font-weight: bold;
                                        color: ${isReturnTrip ? '#e65100' : 'black'}">
                              ${segment.distance} km ${isReturnTrip ? '(Return)' : ''}
                            </div>`,
                      iconSize: [isReturnTrip ? 90 : 60, 20],
                      iconAnchor: [isReturnTrip ? 45 : 30, 10]
                    });
                    
                    return (
                      <Marker 
                        key={`distance-${index}`}
                        position={segment.midpoint}
                        icon={distanceIcon}
                        zIndexOffset={-1000} // Place behind other markers
                      >
                        <Popup>
                          <div>
                            <h4>Segment Information</h4>
                            <p><strong>From:</strong> {segment.startName}</p>
                            <p><strong>To:</strong> {segment.endName}</p>
                            <p><strong>Distance:</strong> {segment.distance} km</p>
                            {isReturnTrip && <p><strong>Type:</strong> Return trip to warehouse</p>}
                          </div>
                        </Popup>
                      </Marker>
                    );
                  })}
                </>
              )}
              
              {/* Display all delivery points */}
              {pointsToDisplay.map((point) => {
                // Find next point in route if this point is in the selected route
                let nextPointInfo = null;
                if (selectedRoute && selectedRoute.delivery_points) {
                  const pointIndex = selectedRoute.delivery_points.indexOf(point.id);
                  if (pointIndex !== -1 && pointIndex < selectedRoute.delivery_points.length - 1) {
                    const nextPointId = selectedRoute.delivery_points[pointIndex + 1];
                    const nextPoint = deliveryPoints.find(p => p.id === nextPointId);
                    
                    if (nextPoint && selectedRoute.segment_distances && 
                        selectedRoute.segment_distances[point.id + '-' + nextPointId]) {
                      nextPointInfo = (
                        <p><strong>Distance to next stop:</strong> {selectedRoute.segment_distances[point.id + '-' + nextPointId]} km</p>
                      );
                    }
                  }
                }
                
                return (
                  <Marker 
                    key={`marker-${point.id}`} 
                    position={[point.latitude, point.longitude]}
                    icon={categoryIcons[point.category]}
                  >
                    <Popup>
                      <div>
                        <strong>{point.description}</strong><br/>
                        {point.address}<br/>
                        {point.pallets > 0 ? `Pallets: ${point.pallets}` : 'Warehouse'}
                        {selectedRoute && selectedRoute.delivery_points && 
                         selectedRoute.delivery_points.includes(point.id) && (
                          <><br/><strong style={{color: 'red'}}>Part of selected route</strong></>
                        )}
                        {nextPointInfo}
                      </div>
                    </Popup>
                  </Marker>
                );
              })}
            </MapContainer>
          </div>
        </Box>
      </Grid>
      
      {/* Route visualization panel - takes 25% of the width */}
      <Grid item xs={12} md={3} sx={{ height: '100vh', overflowY: 'auto', backgroundColor: '#f5f5f5', padding: 2 }}>
        <Paper elevation={3} sx={{ p: 2, height: '100%' }}>
          <Typography variant="h5" gutterBottom sx={{ borderBottom: '2px solid #2196F3', pb: 1 }}>
            Route Visualization
          </Typography>
          
          <FormControl fullWidth sx={{ mb: 3 }}>
            <InputLabel>Select Route</InputLabel>
            <Select
              value={selectedRoute ? selectedRoute.id.toString() : ''}
              onChange={(event) => {
                const routeId = event.target.value;
                if (routeId === '') {
                  setSelectedRoute(null);
                  return;
                }
                
                const route = routes.find(r => r.id.toString() === routeId);
                setSelectedRoute(route);
              }}
              label="Select Route"
            >
              <MenuItem value="">None</MenuItem>
              {routes.map((route) => (
                <MenuItem key={route.id} value={route.id.toString()}>
                  Route #{route.id} - {new Date(route.delivery_date).toLocaleDateString()}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          
          {selectedRoute && (
            <Paper variant="outlined" sx={{ p: 2, mb: 3, bgcolor: '#fff' }}>
              <Typography variant="h6" gutterBottom>
                Route #{selectedRoute.id} Details
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                <Typography variant="body1">
                  <strong>Delivery points:</strong> {selectedRoute.delivery_points ? selectedRoute.delivery_points.length : 0}
                </Typography>
                <Typography variant="body1">
                  <strong>Total pallets:</strong> {selectedRoute.point_pallets ? Object.values(selectedRoute.point_pallets).reduce((sum, val) => sum + val, 0) : 0}
                </Typography>
                <Typography variant="body1">
                  <strong>Total distance:</strong> {selectedRoute.total_distance} km
                </Typography>
                <Typography variant="body1">
                  <strong>Delivery date:</strong> {new Date(selectedRoute.delivery_date).toLocaleDateString()}
                </Typography>
              </Box>
              <Typography variant="caption" sx={{ display: 'block', mt: 1, fontStyle: 'italic', color: 'text.secondary' }}>
                Distance calculated along real roads with return trip to warehouse
              </Typography>
            </Paper>
          )}
          
          {/* Map Legend */}
          <Paper variant="outlined" sx={{ p: 2, bgcolor: '#fff' }}>
            <Typography variant="h6" gutterBottom>
              Map Legend
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 20, height: 4, bgcolor: '#2196F3', mr: 1 }} />
                <Typography variant="body2">Outbound route</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 20, height: 4, bgcolor: '#FF9800', mr: 1, borderTop: '1px dashed white', borderBottom: '1px dashed white' }} />
                <Typography variant="body2">Return trip to warehouse</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 12, height: 12, bgcolor: '#e53935', borderRadius: '50%', mr: 1 }} />
                <Typography variant="body2">Warehouse</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 12, height: 12, bgcolor: '#1976d2', borderRadius: '50%', mr: 1 }} />
                <Typography variant="body2">Blue category (33/18-pallet trucks)</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 12, height: 12, bgcolor: '#2e7d32', borderRadius: '50%', mr: 1 }} />
                <Typography variant="body2">Green category (18-pallet trucks)</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 12, height: 12, bgcolor: '#ed6c02', borderRadius: '50%', mr: 1 }} />
                <Typography variant="body2">Yellow category (15-pallet trucks)</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: 12, height: 12, bgcolor: '#9c27b0', borderRadius: '50%', mr: 1 }} />
                <Typography variant="body2">Purple category (10-pallet trucks)</Typography>
              </Box>
            </Box>
          </Paper>
          
          {selectedRoute && selectedRoute.delivery_points && selectedRoute.delivery_points.length > 0 && (
            <Paper variant="outlined" sx={{ p: 2, mt: 3, bgcolor: '#fff' }}>
              <Typography variant="h6" gutterBottom>
                Delivery Sequence
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                {[1000, ...selectedRoute.delivery_points, 1000].map((pointId, index) => {
                  const point = deliveryPoints.find(p => p.id === pointId);
                  if (!point) return null;
                  
                  return (
                    <Box key={`sequence-${index}`} sx={{ display: 'flex', alignItems: 'center', p: 1, borderLeft: `4px solid ${index === 0 || index === selectedRoute.delivery_points.length + 1 ? '#e53935' : categoryIcons[point.category].options.iconUrl.includes('blue') ? '#1976d2' : categoryIcons[point.category].options.iconUrl.includes('green') ? '#2e7d32' : categoryIcons[point.category].options.iconUrl.includes('yellow') ? '#ed6c02' : '#9c27b0'}` }}>
                      <Typography variant="body2">
                        {index === 0 ? 'Start: ' : index === selectedRoute.delivery_points.length + 1 ? 'End: ' : `${index}. `}
                        <strong>{point.description}</strong>
                        {point.pallets > 0 && ` (${point.pallets} pallets)`}
                      </Typography>
                    </Box>
                  );
                })}
              </Box>
            </Paper>
          )}
        </Paper>
      </Grid>
    </Grid>
  );
};

export default Map;