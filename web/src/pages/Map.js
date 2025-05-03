import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Chip,
  Button,
  Alert,
} from '@mui/material';
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import axios from 'axios';
import L from 'leaflet';

// MapController component to handle map view updates
const MapController = ({ center, zoom }) => {
  const map = useMap();
  
  useEffect(() => {
    map.setView(center, zoom);
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

// Line colors for routes
const routeColors = [
  '#3388ff', // Blue
  '#33a02c', // Green
  '#ff7f00', // Orange
  '#6a3d9a', // Purple
  '#e31a1c', // Red
  '#1f78b4', // Dark Blue
  '#b15928', // Brown
  '#a6cee3', // Light Blue
  '#b2df8a', // Light Green
  '#fb9a99', // Light Red
];

const Map = () => {
  const [deliveryPoints, setDeliveryPoints] = useState([]);
  const [routes, setRoutes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedRoute, setSelectedRoute] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [mapCenter, setMapCenter] = useState([50.45, 30.52]); // Default center (Kyiv)
  const [mapZoom, setMapZoom] = useState(12);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      // Fetch delivery points
      const pointsResponse = await axios.get('/api/delivery-points');
      setDeliveryPoints(pointsResponse.data);
      
      // Fetch routes
      const routesResponse = await axios.get('/api/routes');
      setRoutes(routesResponse.data);
    } catch (error) {
      console.error('Error fetching data:', error);
      setError('Failed to load data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleRouteChange = (event) => {
    setSelectedRoute(event.target.value);
    setSelectedCategory(''); // Clear category filter when route is selected
    
    // If a route is selected, center the map on the first point of the route
    if (event.target.value) {
      const route = routes.find(r => r.id.toString() === event.target.value);
      if (route && route.delivery_points.length > 0) {
        const firstPointId = route.delivery_points[0];
        const point = deliveryPoints.find(p => p.id === firstPointId);
        if (point && point.latitude && point.longitude) {
          setMapCenter([point.latitude, point.longitude]);
          setMapZoom(13); // Zoom in a bit
        }
      }
    } else {
      // Reset to default view
      setMapCenter([50.45, 30.52]);
      setMapZoom(12);
    }
  };

  const handleCategoryChange = (event) => {
    setSelectedCategory(event.target.value);
    setSelectedRoute(''); // Clear route filter when category is selected
    
    // If a category is selected, adjust the map view to show all points of that category
    if (event.target.value) {
      const categoryPoints = deliveryPoints.filter(p => p.category === event.target.value);
      if (categoryPoints.length > 0) {
        // Find the average lat/lng to center the map
        const avgLat = categoryPoints.reduce((sum, p) => sum + p.latitude, 0) / categoryPoints.length;
        const avgLng = categoryPoints.reduce((sum, p) => sum + p.longitude, 0) / categoryPoints.length;
        setMapCenter([avgLat, avgLng]);
        setMapZoom(12); // Slightly zoomed out to see all points
      }
    } else {
      // Reset to default view
      setMapCenter([50.45, 30.52]);
      setMapZoom(12);
    }
  };

  const getFilteredPoints = () => {
    if (selectedRoute) {
      const route = routes.find(r => r.id.toString() === selectedRoute);
      if (route) {
        return deliveryPoints.filter(p => route.delivery_points.includes(p.id));
      }
      return [];
    }
    
    if (selectedCategory) {
      return deliveryPoints.filter(p => p.category === selectedCategory);
    }
    
    return deliveryPoints;
  };

  const getRoutePolylines = () => {
    if (!selectedRoute) return [];
    
    const route = routes.find(r => r.id.toString() === selectedRoute);
    if (!route || !route.delivery_points || route.delivery_points.length < 2) return [];
    
    // Get all points in the route with their coordinates
    const routePoints = route.delivery_points.map(pointId => {
      const point = deliveryPoints.find(p => p.id === pointId);
      return point ? [point.latitude, point.longitude] : null;
    }).filter(point => point !== null);
    
    // Create a polyline for the route
    return [{
      positions: routePoints,
      color: routeColors[route.id % routeColors.length],
      id: route.id,
      weight: 4,
      opacity: 0.8,
      dashArray: '10, 5' // Create a dashed line for better visibility
    }];
  };

  const getCategoryColor = (category) => {
    switch (category) {
      case 'Blue':
        return 'primary';
      case 'Green':
        return 'success';
      case 'Yellow':
        return 'warning';
      case 'Purple':
        return 'secondary';
      default:
        return 'default';
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Delivery Map
      </Typography>
      
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={4}>
            <FormControl fullWidth>
              <InputLabel id="route-select-label">Filter by Route</InputLabel>
              <Select
                labelId="route-select-label"
                id="route-select"
                value={selectedRoute}
                label="Filter by Route"
                onChange={handleRouteChange}
              >
                <MenuItem value="">All Routes</MenuItem>
                {routes.map((route) => (
                  <MenuItem key={route.id} value={route.id.toString()}>
                    Route {route.id} ({new Date(route.delivery_date).toLocaleDateString()})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <FormControl fullWidth>
              <InputLabel id="category-select-label">Filter by Category</InputLabel>
              <Select
                labelId="category-select-label"
                id="category-select"
                value={selectedCategory}
                label="Filter by Category"
                onChange={handleCategoryChange}
                disabled={!!selectedRoute}
              >
                <MenuItem value="">All Categories</MenuItem>
                <MenuItem value="Blue">
                  <Chip label="Blue" color="primary" size="small" sx={{ mr: 1 }} /> Blue (33-pallet)
                </MenuItem>
                <MenuItem value="Green">
                  <Chip label="Green" color="success" size="small" sx={{ mr: 1 }} /> Green (18-pallet)
                </MenuItem>
                <MenuItem value="Yellow">
                  <Chip label="Yellow" color="warning" size="small" sx={{ mr: 1 }} /> Yellow (15-pallet)
                </MenuItem>
                <MenuItem value="Purple">
                  <Chip label="Purple" color="secondary" size="small" sx={{ mr: 1 }} /> Purple (10-pallet)
                </MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Button 
              variant="outlined" 
              onClick={fetchData} 
              fullWidth
            >
              Refresh Data
            </Button>
          </Grid>
        </Grid>
      </Paper>
      
      <Paper sx={{ height: '70vh', width: '100%', overflow: 'hidden' }}>
        <MapContainer 
          center={mapCenter} 
          zoom={mapZoom} 
          style={{ height: '100%', width: '100%' }}
          key={`${mapCenter[0]}-${mapCenter[1]}-${mapZoom}`} // Add key to force re-render when center or zoom changes
        >
          <MapController center={mapCenter} zoom={mapZoom} />
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          
          {getFilteredPoints().map((point) => (
            <Marker 
              key={point.id} 
              position={[point.latitude, point.longitude]}
              icon={categoryIcons[point.category]}
            >
              <Popup>
                <div>
                  <Typography variant="subtitle1">
                    <strong>{point.description}</strong>
                  </Typography>
                  <Typography variant="body2">
                    Category: <Chip 
                      label={point.category} 
                      color={getCategoryColor(point.category)} 
                      size="small" 
                    />
                  </Typography>
                  <Typography variant="body2">Address: {point.address}</Typography>
                  <Typography variant="body2">Pallets: {point.pallets}</Typography>
                  <Typography variant="body2">
                    Coordinates: {point.latitude.toFixed(5)}, {point.longitude.toFixed(5)}
                  </Typography>
                </div>
              </Popup>
            </Marker>
          ))}
          
          {getRoutePolylines().map((line) => (
            <Polyline 
              key={`route-${line.id}`}
              positions={line.positions} 
              color={line.color} 
              weight={line.weight || 4}
              opacity={line.opacity || 0.7}
              dashArray={line.dashArray}
            />
          ))}
        </MapContainer>
      </Paper>
    </Box>
  );
};

export default Map;