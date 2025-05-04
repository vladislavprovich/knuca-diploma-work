import React, { useState, useEffect } from 'react';
import { Box } from '@mui/material';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import './Map.css';
import L from 'leaflet';

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
  const [mapCenter] = useState([50.3833, 30.3667]); // Default center (Vyshneve warehouse)
  const [mapZoom] = useState(10);
  
  // Force map to load properly by using useEffect
  useEffect(() => {
    // This helps ensure the map container is properly sized when the component mounts
    const mapContainer = document.querySelector('.leaflet-container');
    if (mapContainer) {
      window.dispatchEvent(new Event('resize'));
    }
  }, []);

  return (
    <Box sx={mapWrapperStyle} className="map-wrapper">
      {/* Wrap MapContainer in a div to ensure it has a parent React component context */}
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
          
          {staticDeliveryPoints.map((point) => (
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
                </div>
              </Popup>
            </Marker>
          ))}
        </MapContainer>
      </div>
    </Box>
  );
};

export default Map;