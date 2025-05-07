import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Chip,
  CircularProgress,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  Checkbox,
} from '@mui/material';
import axios from 'axios';

const Trucks = () => {
  const [trucks, setTrucks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('');
  const [capacityFilter, setCapacityFilter] = useState('');
  const [availabilityFilter, setAvailabilityFilter] = useState('');
  const [addTruckDialogOpen, setAddTruckDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [currentTruck, setCurrentTruck] = useState(null);
  const [deleting, setDeleting] = useState(false);
  const [newTruck, setNewTruck] = useState({
    license_plate: '',
    brand: '',
    capacity: 10,
    driver_name: '',
    driver_surname: '',
    has_trailer: false,
    trailer_info: '',
  });
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchTrucks();
  }, []);

  const fetchTrucks = async () => {
    setLoading(true);
    try {
      const response = await axios.get('/api/trucks');
      setTrucks(response.data);
    } catch (error) {
      console.error('Error fetching trucks:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddTruckClick = () => {
    setAddTruckDialogOpen(true);
  };

  const handleCloseAddTruck = () => {
    setAddTruckDialogOpen(false);
    // Reset form
    setNewTruck({
      license_plate: '',
      brand: '',
      capacity: 10,
      driver_name: '',
      driver_surname: '',
      has_trailer: false,
      trailer_info: '',
    });
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewTruck({
      ...newTruck,
      [name]: value,
    });
  };

  const handleCheckboxChange = (e) => {
    const { name, checked } = e.target;
    setNewTruck({
      ...newTruck,
      [name]: checked,
    });
  };

  const handleCapacityChange = (e) => {
    const capacity = parseInt(e.target.value);
    const hasTrailer = capacity === 33; // Only 33-pallet trucks have trailers
    
    setNewTruck({
      ...newTruck,
      capacity,
      has_trailer: hasTrailer,
      trailer_info: hasTrailer ? newTruck.trailer_info : '',
    });
  };

  const handleSubmit = async () => {
    // Validate form
    if (!newTruck.license_plate || !newTruck.brand || !newTruck.driver_name || !newTruck.driver_surname) {
      alert('Please fill in all required fields');
      return;
    }

    setSubmitting(true);
    try {
      await axios.post('/api/trucks', newTruck);
      handleCloseAddTruck();
      fetchTrucks(); // Refresh the truck list
    } catch (error) {
      console.error('Error creating truck:', error);
      alert('Failed to create truck. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleDeleteClick = (truck) => {
    setCurrentTruck(truck);
    setDeleteDialogOpen(true);
  };

  const handleCloseDeleteDialog = () => {
    setDeleteDialogOpen(false);
    setCurrentTruck(null);
  };

  const handleDeleteTruck = async () => {
    if (!currentTruck) return;
    
    setDeleting(true);
    try {
      await axios.delete(`/api/trucks/${currentTruck.id}`);
      fetchTrucks(); // Refresh the truck list
      setDeleteDialogOpen(false);
      setCurrentTruck(null);
    } catch (error) {
      console.error('Error deleting truck:', error);
      alert('Failed to delete truck. Please try again.');
    } finally {
      setDeleting(false);
    }
  };

  const filteredTrucks = trucks.filter((truck) => {
    const matchesSearch = filter === '' || 
      truck.id.toString().includes(filter) || 
      truck.license_plate.toLowerCase().includes(filter.toLowerCase()) ||
      truck.driver_name.toLowerCase().includes(filter.toLowerCase()) ||
      truck.driver_surname.toLowerCase().includes(filter.toLowerCase()) ||
      truck.brand.toLowerCase().includes(filter.toLowerCase());
    
    const matchesCapacity = capacityFilter === '' || truck.capacity.toString() === capacityFilter;
    const matchesAvailability = availabilityFilter === '' || 
      (availabilityFilter === 'available' && truck.available) ||
      (availabilityFilter === 'unavailable' && !truck.available);
    
    return matchesSearch && matchesCapacity && matchesAvailability;
  });

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Trucks Fleet
      </Typography>
      
      {/* Action Buttons */}
      <Box sx={{ mb: 3 }}>
        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleAddTruckClick}
        >
          Add New Truck
        </Button>
      </Box>
      
      {/* Filters */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={4}>
          <TextField
            fullWidth
            label="Search"
            variant="outlined"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="Search by ID, license plate, driver, or brand"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <FormControl fullWidth>
            <InputLabel>Capacity</InputLabel>
            <Select
              value={capacityFilter}
              label="Capacity"
              onChange={(e) => setCapacityFilter(e.target.value)}
            >
              <MenuItem value="">All Capacities</MenuItem>
              <MenuItem value="33">33 Pallets</MenuItem>
              <MenuItem value="18">18 Pallets</MenuItem>
              <MenuItem value="15">15 Pallets</MenuItem>
              <MenuItem value="10">10 Pallets</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <FormControl fullWidth>
            <InputLabel>Availability</InputLabel>
            <Select
              value={availabilityFilter}
              label="Availability"
              onChange={(e) => setAvailabilityFilter(e.target.value)}
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="available">Available</MenuItem>
              <MenuItem value="unavailable">Unavailable</MenuItem>
            </Select>
          </FormControl>
        </Grid>
      </Grid>
      
      {/* Trucks Table */}
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>License Plate</TableCell>
              <TableCell>Brand</TableCell>
              <TableCell>Capacity</TableCell>
              <TableCell>Driver</TableCell>
              <TableCell>Trailer</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredTrucks.map((truck) => (
              <TableRow key={truck.id}>
                <TableCell>{truck.id}</TableCell>
                <TableCell>{truck.license_plate}</TableCell>
                <TableCell>{truck.brand}</TableCell>
                <TableCell>{truck.capacity} pallets</TableCell>
                <TableCell>{`${truck.driver_name} ${truck.driver_surname}`}</TableCell>
                <TableCell>
                  {truck.has_trailer ? (
                    <Chip 
                      label={truck.trailer_info || 'Has Trailer'} 
                      color="info" 
                      size="small" 
                    />
                  ) : (
                    'No Trailer'
                  )}
                </TableCell>
                <TableCell>
                  <Chip 
                    label={truck.available ? 'Available' : 'In Use'} 
                    color={truck.available ? 'success' : 'error'} 
                    size="small" 
                  />
                </TableCell>
                <TableCell>
                  <Button
                    variant="contained"
                    color="error"
                    size="small"
                    onClick={() => handleDeleteClick(truck)}
                  >
                    Delete
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Delete Truck Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={handleCloseDeleteDialog}>
        <DialogTitle>Delete Truck</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete the truck with driver {currentTruck?.driver_name} {currentTruck?.driver_surname}? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog}>Cancel</Button>
          <Button 
            onClick={handleDeleteTruck} 
            variant="contained" 
            color="error"
            disabled={deleting}
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Truck Dialog */}
      <Dialog open={addTruckDialogOpen} onClose={handleCloseAddTruck} maxWidth="sm" fullWidth>
        <DialogTitle>Add New Truck</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="License Plate"
                name="license_plate"
                value={newTruck.license_plate}
                onChange={handleInputChange}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Brand"
                name="brand"
                value={newTruck.brand}
                onChange={handleInputChange}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Driver Name"
                name="driver_name"
                value={newTruck.driver_name}
                onChange={handleInputChange}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Driver Surname"
                name="driver_surname"
                value={newTruck.driver_surname}
                onChange={handleInputChange}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Capacity</InputLabel>
                <Select
                  value={newTruck.capacity}
                  label="Capacity"
                  onChange={handleCapacityChange}
                >
                  <MenuItem value={10}>10 Pallets</MenuItem>
                  <MenuItem value={15}>15 Pallets</MenuItem>
                  <MenuItem value={18}>18 Pallets</MenuItem>
                  <MenuItem value={33}>33 Pallets</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            {newTruck.capacity === 33 && (
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Trailer Info"
                  name="trailer_info"
                  value={newTruck.trailer_info}
                  onChange={handleInputChange}
                  placeholder="Optional"
                />
              </Grid>
            )}
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseAddTruck}>Cancel</Button>
          <Button 
            onClick={handleSubmit} 
            variant="contained" 
            color="primary"
            disabled={submitting}
          >
            {submitting ? 'Adding...' : 'Add Truck'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Trucks;