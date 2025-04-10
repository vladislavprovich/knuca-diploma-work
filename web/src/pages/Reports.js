import React, { useState, useEffect } from 'react';
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
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import dayjs from 'dayjs';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js';
import { Bar } from 'react-chartjs-2';
import axios from 'axios';

// Register ChartJS components
ChartJS.register(ArcElement, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

const Reports = () => {
  const [loading, setLoading] = useState(true);
  const [reportType, setReportType] = useState('daily');
  const [startDate, setStartDate] = useState(dayjs().subtract(7, 'day'));
  const [endDate, setEndDate] = useState(dayjs());
  const [reportData, setReportData] = useState(null);
  const [trucks, setTrucks] = useState([]);
  const [deliveryPoints, setDeliveryPoints] = useState([]);
  const [routes, setRoutes] = useState([]);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      // Fetch all necessary data
      const [trucksResponse, pointsResponse, routesResponse] = await Promise.all([
        axios.get('/api/trucks'),
        axios.get('/api/delivery-points'),
        axios.get('/api/routes')
      ]);
      
      setTrucks(trucksResponse.data);
      setDeliveryPoints(pointsResponse.data);
      setRoutes(routesResponse.data);
      
      // Generate initial report
      generateReport(reportType, startDate, endDate, routesResponse.data);
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateReport = () => {
    generateReport(reportType, startDate, endDate, routes);
  };

  const generateReport = (type, start, end, routesData) => {
    // Filter routes by date range
    const filteredRoutes = routesData.filter(route => {
      const routeDate = dayjs(route.delivery_date);
      return routeDate.isAfter(start) && routeDate.isBefore(end.add(1, 'day'));
    });
    
    let reportResult = {};
    
    switch (type) {
      case 'daily':
        // Group routes by day
        reportResult = filteredRoutes.reduce((acc, route) => {
          const day = dayjs(route.delivery_date).format('YYYY-MM-DD');
          if (!acc[day]) {
            acc[day] = {
              routes: 0,
              totalDistance: 0,
              totalCost: 0,
              deliveryPoints: 0
            };
          }
          
          acc[day].routes += 1;
          acc[day].totalDistance += route.total_distance;
          acc[day].totalCost += route.total_cost;
          acc[day].deliveryPoints += route.delivery_points.length;
          
          return acc;
        }, {});
        break;
        
      case 'truck':
        // Group by truck
        reportResult = filteredRoutes.reduce((acc, route) => {
          const truckId = route.truck_id;
          const truck = trucks.find(t => t.id === truckId);
          const truckName = truck ? `${truck.brand} (${truck.license_plate})` : `Truck ${truckId}`;
          
          if (!acc[truckName]) {
            acc[truckName] = {
              routes: 0,
              totalDistance: 0,
              totalCost: 0,
              deliveryPoints: 0
            };
          }
          
          acc[truckName].routes += 1;
          acc[truckName].totalDistance += route.total_distance;
          acc[truckName].totalCost += route.total_cost;
          acc[truckName].deliveryPoints += route.delivery_points.length;
          
          return acc;
        }, {});
        break;
        
      case 'category':
        // Group by delivery point category
        reportResult = {};
        
        // Initialize categories
        ['Blue', 'Green', 'Yellow', 'Purple'].forEach(category => {
          reportResult[category] = {
            points: 0,
            pallets: 0
          };
        });
        
        // Count points and pallets by category
        deliveryPoints.forEach(point => {
          if (point.pallets > 0 && reportResult[point.category]) {
            reportResult[point.category].points += 1;
            reportResult[point.category].pallets += point.pallets;
          }
        });
        break;
        
      default:
        break;
    }
    
    setReportData(reportResult);
  };

  const renderReportTable = () => {
    if (!reportData || Object.keys(reportData).length === 0) {
      return (
        <Typography variant="body1" sx={{ mt: 2 }}>
          No data available for the selected period.
        </Typography>
      );
    }
    
    switch (reportType) {
      case 'daily':
        return (
          <TableContainer component={Paper} sx={{ mt: 3 }}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Date</TableCell>
                  <TableCell>Routes</TableCell>
                  <TableCell>Delivery Points</TableCell>
                  <TableCell>Total Distance (km)</TableCell>
                  <TableCell>Total Cost</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {Object.entries(reportData)
                  .sort(([dateA], [dateB]) => dayjs(dateA).diff(dayjs(dateB)))
                  .map(([date, data]) => (
                    <TableRow key={date}>
                      <TableCell>{date}</TableCell>
                      <TableCell>{data.routes}</TableCell>
                      <TableCell>{data.deliveryPoints}</TableCell>
                      <TableCell>{data.totalDistance}</TableCell>
                      <TableCell>{data.totalCost}</TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
        );
        
      case 'truck':
        return (
          <TableContainer component={Paper} sx={{ mt: 3 }}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Truck</TableCell>
                  <TableCell>Routes</TableCell>
                  <TableCell>Delivery Points</TableCell>
                  <TableCell>Total Distance (km)</TableCell>
                  <TableCell>Total Cost</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {Object.entries(reportData)
                  .sort(([, dataA], [, dataB]) => dataB.totalDistance - dataA.totalDistance)
                  .map(([truck, data]) => (
                    <TableRow key={truck}>
                      <TableCell>{truck}</TableCell>
                      <TableCell>{data.routes}</TableCell>
                      <TableCell>{data.deliveryPoints}</TableCell>
                      <TableCell>{data.totalDistance}</TableCell>
                      <TableCell>{data.totalCost}</TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
        );
        
      case 'category':
        return (
          <TableContainer component={Paper} sx={{ mt: 3 }}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Category</TableCell>
                  <TableCell>Delivery Points</TableCell>
                  <TableCell>Total Pallets</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {Object.entries(reportData)
                  .sort(([, dataA], [, dataB]) => dataB.pallets - dataA.pallets)
                  .map(([category, data]) => (
                    <TableRow key={category}>
                      <TableCell>{category}</TableCell>
                      <TableCell>{data.points}</TableCell>
                      <TableCell>{data.pallets}</TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
        );
        
      default:
        return null;
    }
  };

  const renderChart = () => {
    if (!reportData || Object.keys(reportData).length === 0) {
      return null;
    }
    
    let chartData = {};
    
    switch (reportType) {
      case 'daily':
        chartData = {
          labels: Object.keys(reportData).sort((a, b) => dayjs(a).diff(dayjs(b))),
          datasets: [
            {
              label: 'Total Distance (km)',
              data: Object.entries(reportData)
                .sort(([dateA], [dateB]) => dayjs(dateA).diff(dayjs(dateB)))
                .map(([, data]) => data.totalDistance),
              backgroundColor: 'rgba(54, 162, 235, 0.6)',
              borderColor: 'rgba(54, 162, 235, 1)',
              borderWidth: 1,
            },
            {
              label: 'Total Cost',
              data: Object.entries(reportData)
                .sort(([dateA], [dateB]) => dayjs(dateA).diff(dayjs(dateB)))
                .map(([, data]) => data.totalCost),
              backgroundColor: 'rgba(255, 99, 132, 0.6)',
              borderColor: 'rgba(255, 99, 132, 1)',
              borderWidth: 1,
            },
          ],
        };
        break;
        
      case 'truck':
        chartData = {
          labels: Object.keys(reportData),
          datasets: [
            {
              label: 'Total Distance (km)',
              data: Object.values(reportData).map(data => data.totalDistance),
              backgroundColor: 'rgba(54, 162, 235, 0.6)',
              borderColor: 'rgba(54, 162, 235, 1)',
              borderWidth: 1,
            },
            {
              label: 'Total Cost',
              data: Object.values(reportData).map(data => data.totalCost),
              backgroundColor: 'rgba(255, 99, 132, 0.6)',
              borderColor: 'rgba(255, 99, 132, 1)',
              borderWidth: 1,
            },
          ],
        };
        break;
        
      case 'category':
        chartData = {
          labels: Object.keys(reportData),
          datasets: [
            {
              label: 'Total Pallets',
              data: Object.values(reportData).map(data => data.pallets),
              backgroundColor: [
                'rgba(54, 162, 235, 0.6)',
                'rgba(75, 192, 192, 0.6)',
                'rgba(255, 206, 86, 0.6)',
                'rgba(153, 102, 255, 0.6)',
              ],
              borderColor: [
                'rgba(54, 162, 235, 1)',
                'rgba(75, 192, 192, 1)',
                'rgba(255, 206, 86, 1)',
                'rgba(153, 102, 255, 1)',
              ],
              borderWidth: 1,
            },
          ],
        };
        break;
        
      default:
        return null;
    }
    
    return (
      <Box sx={{ height: 400, mt: 3 }}>
        <Bar 
          data={chartData} 
          options={{
            responsive: true,
            maintainAspectRatio: false,
            scales: {
              y: {
                beginAtZero: true
              }
            }
          }} 
        />
      </Box>
    );
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
        Reports
      </Typography>
      
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          Generate Report
        </Typography>
        
        <Grid container spacing={2}>
          <Grid item xs={12} sm={4}>
            <FormControl fullWidth margin="normal">
              <InputLabel>Report Type</InputLabel>
              <Select
                value={reportType}
                label="Report Type"
                onChange={(e) => setReportType(e.target.value)}
              >
                <MenuItem value="daily">Daily Report</MenuItem>
                <MenuItem value="truck">Truck Performance</MenuItem>
                <MenuItem value="category">Category Analysis</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} sm={4}>
            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DatePicker
                label="Start Date"
                value={startDate}
                onChange={(newValue) => setStartDate(newValue)}
                slotProps={{
                  textField: {
                    fullWidth: true,
                    margin: "normal"
                  }
                }}
              />
            </LocalizationProvider>
          </Grid>
          
          <Grid item xs={12} sm={4}>
            <LocalizationProvider dateAdapter={AdapterDayjs}>
              <DatePicker
                label="End Date"
                value={endDate}
                onChange={(newValue) => setEndDate(newValue)}
                slotProps={{
                  textField: {
                    fullWidth: true,
                    margin: "normal"
                  }
                }}
              />
            </LocalizationProvider>
          </Grid>
          
          <Grid item xs={12}>
            <Button 
              variant="contained" 
              color="primary" 
              onClick={handleGenerateReport}
            >
              Generate Report
            </Button>
          </Grid>
        </Grid>
      </Paper>
      
      {reportData && (
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            {reportType === 'daily' && 'Daily Report'}
            {reportType === 'truck' && 'Truck Performance Report'}
            {reportType === 'category' && 'Category Analysis Report'}
          </Typography>
          
          {renderReportTable()}
          {renderChart()}
        </Paper>
      )}
    </Box>
  );
};

export default Reports;