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
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Grid,
  Chip,
} from '@mui/material';
import axios from 'axios';

const DeliveryPoints = () => {
  const [deliveryPoints, setDeliveryPoints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('');
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [currentPoint, setCurrentPoint] = useState(null);
  const [pallets, setPallets] = useState(0);

  useEffect(() => {
    fetchDeliveryPoints();
  }, []);

  const fetchDeliveryPoints = async () => {
    setLoading(true);
    try {
      const response = await axios.get('/api/delivery-points');
      setDeliveryPoints(response.data);
    } catch (error) {
      console.error('Error fetching delivery points:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleEditClick = (point) => {
    setCurrentPoint(point);
    setPallets(point.pallets);
    setEditDialogOpen(true);
  };

  const handleSave = async () => {
    try {
      await axios.put(`/api/delivery-points/${currentPoint.id}/pallets`, {
        pallets: parseInt(pallets, 10),
      });
      setEditDialogOpen(false);
      fetchDeliveryPoints(); // Refresh the list
    } catch (error) {
      console.error('Error updating delivery point:', error);
    }
  };

  const handleClose = () => {
    setEditDialogOpen(false);
  };

  const filteredPoints = deliveryPoints.filter((point) => {
    const matchesSearch = filter === '' || 
      point.id.toString().includes(filter) || 
      (point.description && point.description.toLowerCase().includes(filter.toLowerCase())) ||
      (point.address && point.address.toLowerCase().includes(filter.toLowerCase()));
    
    const matchesCategory = categoryFilter === '' || point.category === categoryFilter;
    
    return matchesSearch && matchesCategory;
  });

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
    <Box>
      <Typography variant="h4" gutterBottom>
        Delivery Points
      </Typography>
      
      {/* Filters */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={4}>
          <TextField
            fullWidth
            label="Search"
            variant="outlined"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="Search by ID, description, or address"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <FormControl fullWidth>
            <InputLabel>Category</InputLabel>
            <Select
              value={categoryFilter}
              label="Category"
              onChange={(e) => setCategoryFilter(e.target.value)}
            >
              <MenuItem value="">All Categories</MenuItem>
              <MenuItem value="Blue">Blue</MenuItem>
              <MenuItem value="Green">Green</MenuItem>
              <MenuItem value="Yellow">Yellow</MenuItem>
              <MenuItem value="Purple">Purple</MenuItem>
            </Select>
          </FormControl>
        </Grid>
      </Grid>
      
      {/* Delivery Points Table */}
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Category</TableCell>
              <TableCell>Pallets</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Address</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredPoints.map((point) => (
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
                <TableCell>{point.address || '-'}</TableCell>
                <TableCell>
                  <Button 
                    variant="contained" 
                    size="small" 
                    onClick={() => handleEditClick(point)}
                  >
                    Edit Pallets
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      
      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onClose={handleClose}>
        <DialogTitle>Edit Delivery Point</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2 }}>
            <Typography variant="body1" gutterBottom>
              ID: {currentPoint?.id}
            </Typography>
            <Typography variant="body1" gutterBottom>
              Category: {currentPoint?.category}
            </Typography>
            <TextField
              autoFocus
              margin="dense"
              label="Pallets"
              type="number"
              fullWidth
              variant="outlined"
              value={pallets}
              onChange={(e) => setPallets(e.target.value)}
              InputProps={{ inputProps: { min: 0 } }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleSave} variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default DeliveryPoints;