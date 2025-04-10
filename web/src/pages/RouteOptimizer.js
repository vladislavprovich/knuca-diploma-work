import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Alert,
  Snackbar,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import dayjs from 'dayjs';
import axios from 'axios';

const RouteOptimizer = () => {
  const [deliveryPoints, setDeliveryPoints] = useState([]);
  const [trucks, setTrucks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [selectedPoint, setSelectedPoint] = useState('');
  const [pallets, setPallets] = useState(0);
  const [deliveryDate, setDeliveryDate] = useState(dayjs().add(1, 'day'));
  const [optimizing, setOptimizing] = useState(false);
  const [optimizationResult, setOptimizationResult] = useState(null);
  const [resultDialogOpen, setResultDialogOpen] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      // Fetch delivery points
      const pointsResponse = await axios.get('/api/delivery-points');
      setDeliveryPoints(pointsResponse.data);
      
      // Fetch trucks
      const trucksResponse = await axios.get('/api/trucks');
      setTrucks(trucksResponse.data);
    } catch (error) {
      console.error('Error fetching data:', error);
      setSnackbar({
        open: true,
        message: 'Error loading data. Please try again.',
        severity: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCategoryChange = (event) => {
    setSelectedCategory(event.target.value);
    setSelectedPoint('');
  };

  const handlePointChange = (event) => {
    const pointId = event.target.value;
    setSelectedPoint(pointId);
    
    // Set current pallets value if point is selected
    if (pointId) {
      const point = deliveryPoints.find(p => p.id.toString() === pointId);
      if (point) {
        setPallets(point.pallets);
      }
    } else {
      setPallets(0);
    }
  };

  const handleUpdatePallets = async () => {
    if (!selectedPoint || pallets < 0) {
      setSnackbar({
        open: true,
        message: 'Please select a delivery point and enter a valid number of pallets.',
        severity: 'error'
      });
      return;
    }

    try {
      await axios.put(`/api/delivery-points/${selectedPoint}/pallets`, {
        pallets: parseInt(pallets, 10),
      });
      
      // Update local state
      setDeliveryPoints(prevPoints => 
        prevPoints.map(point => 
          point.id.toString() === selectedPoint 
            ? { ...point, pallets: parseInt(pallets, 10) } 
            : point
        )
      );
      
      setSnackbar({
        open: true,
        message: 'Pallets updated successfully!',
        severity: 'success'
      });
    } catch (error) {
      console.error('Error updating pallets:', error);
      setSnackbar({
        open: true,
        message: 'Error updating pallets. Please try again.',
        severity: 'error'
      });
    }
  };

  const handleOptimizeRoutes = async () => {
    // Check if there are any delivery points with pallets
    const hasPointsWithPallets = deliveryPoints.some(point => point.pallets > 0);
    
    if (!hasPointsWithPallets) {
      setSnackbar({
        open: true,
        message: 'No delivery points have pallets assigned. Please assign pallets before optimizing routes.',
        severity: 'warning'
      });
      return;
    }

    setOptimizing(true);
    try {
      const response = await axios.post('/api/routes/optimize', {
        delivery_date: deliveryDate.format('YYYY-MM-DD')
      });
      
      setOptimizationResult(response.data);
      setResultDialogOpen(true);
      
      // Refresh delivery points to get updated pallets
      fetchData();
    } catch (error) {
      console.error('Error optimizing routes:', error);
      setSnackbar({
        open: true,
        message: 'Error optimizing routes. Please try again.',
        severity: 'error'
      });
    } finally {
      setOptimizing(false);
    }
  };

  const handleCloseResultDialog = () => {
    setResultDialogOpen(false);
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  const filteredPoints = selectedCategory
    ? deliveryPoints.filter(point => point.category === selectedCategory)
    : deliveryPoints;

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

  const getTruckById = (id) => {
    return trucks.find(truck => truck.id === id) || {};
  };

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
        Route Optimizer
      </Typography>
      
      <Grid container spacing={3}>
        {/* Left Panel - Delivery Point Selection */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Assign Pallets to Delivery Points
            </Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <FormControl fullWidth margin="normal">
                  <InputLabel>Category</InputLabel>
                  <Select
                    value={selectedCategory}
                    label="Category"
                    onChange={handleCategoryChange}
                  >
                    <MenuItem value="">All Categories</MenuItem>
                    <MenuItem value="Blue">Blue</MenuItem>
                    <MenuItem value="Green">Green</MenuItem>
                    <MenuItem value="Yellow">Yellow</MenuItem>
                    <MenuItem value="Purple">Purple</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <FormControl fullWidth margin="normal">
                  <InputLabel>Delivery Point</InputLabel>
                  <Select
                    value={selectedPoint}
                    label="Delivery Point"
                    onChange={handlePointChange}
                  >
                    <MenuItem value="">Select a Point</MenuItem>
                    {filteredPoints.map((point) => (
                      <MenuItem key={point.id} value={point.id.toString()}>
                        {point.id} - {point.category} {point.description ? `(${point.description})` : ''}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Pallets"
                  type="number"
                  margin="normal"
                  value={pallets}
                  onChange={(e) => setPallets(e.target.value)}
                  InputProps={{ inputProps: { min: 0 } }}
                  disabled={!selectedPoint}
                />
              </Grid>
              
              <Grid item xs={12} sm={6} sx={{ display: 'flex', alignItems: 'center' }}>
                <Button
                  variant="contained"
                  color="primary"
                  onClick={handleUpdatePallets}
                  disabled={!selectedPoint}
                  sx={{ mt: 2 }}
                  fullWidth
                >
                  Update Pallets
                </Button>
              </Grid>
            </Grid>
          </Paper>
          
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Generate Optimized Routes
            </Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                  <DatePicker
                    label="Delivery Date"
                    value={deliveryDate}
                    onChange={(newValue) => setDeliveryDate(newValue)}
                    slotProps={{
                      textField: {
                        fullWidth: true,
                        margin: "normal"
                      }
                    }}
                  />
                </LocalizationProvider>
              </Grid>
              
              <Grid item xs={12} sm={6} sx={{ display: 'flex', alignItems: 'center' }}>
                <Button
                  variant="contained"
                  color="secondary"
                  onClick={handleOptimizeRoutes}
                  disabled={optimizing}
                  sx={{ mt: 2 }}
                  fullWidth
                >
                  {optimizing ? 'Optimizing...' : 'Optimize Routes'}
                </Button>
              </Grid>
            </Grid>
          </Paper>
        </Grid>
        
        {/* Right Panel - Delivery Points Table */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Delivery Points with Pallets
            </Typography>
            
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>ID</TableCell>
                    <TableCell>Category</TableCell>
                    <TableCell>Pallets</TableCell>
                    <TableCell>Description</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {deliveryPoints
                    .filter(point => point.pallets > 0)
                    .sort((a, b) => b.pallets - a.pallets)
                    .map((point) => (
                      <TableRow key={point.id}>
                        <TableCell>{point.id}</TableCell>
                        <TableCell>
                          <Chip 
                            label={point.category} 
                            color={getCategoryColor(point.category)} 
                            size="small" 
                          />
                        </TableCell>
                        <TableCell>{point.pallets}</TableCell>
                        <TableCell>{point.description || '-'}</TableCell>
                      </TableRow>
                    ))}
                  {deliveryPoints.filter(point => point.pallets > 0).length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} align="center">
                        No delivery points with pallets assigned
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          </Paper>
        </Grid>
      </Grid>
      
      {/* Optimization Result Dialog */}
      <Dialog open={resultDialogOpen} onClose={handleCloseResultDialog} maxWidth="md" fullWidth>
        <DialogTitle>Route Optimization Results</DialogTitle>
        <DialogContent>
          {optimizationResult && optimizationResult.length > 0 ? (
            <Box sx={{ pt: 2 }}>
              <Typography variant="body1" gutterBottom>
                Successfully created {optimizationResult.length} optimized routes for {deliveryDate.format('YYYY-MM-DD')}.
              </Typography>
              
              <TableContainer component={Paper} sx={{ mt: 2 }}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Route ID</TableCell>
                      <TableCell>Truck</TableCell>
                      <TableCell>Points</TableCell>
                      <TableCell>Distance</TableCell>
                      <TableCell>Cost</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {optimizationResult.map((route) => {
                      const truck = getTruckById(route.truck_id);
                      return (
                        <TableRow key={route.id}>
                          <TableCell>{route.id}</TableCell>
                          <TableCell>
                            {truck.license_plate} ({truck.capacity} pallets)
                          </TableCell>
                          <TableCell>{route.delivery_points.length}</TableCell>
                          <TableCell>{route.total_distance} km</TableCell>
                          <TableCell>{route.total_cost} units</TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          ) : (
            <Typography variant="body1">
              No routes were created. Please ensure there are delivery points with pallets assigned.
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseResultDialog}>Close</Button>
        </DialogActions>
      </Dialog>
      
      {/* Snackbar for notifications */}
      <Snackbar 
        open={snackbar.open} 
        autoHideDuration={6000} 
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default RouteOptimizer;