import React, { useState, useEffect } from 'react';
import { Grid, Paper, Typography, Box, Card, CardContent, CircularProgress } from '@mui/material';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js';
import { Pie, Bar } from 'react-chartjs-2';
import axios from 'axios';

// Register ChartJS components
ChartJS.register(ArcElement, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

const Dashboard = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalDeliveryPoints: 0,
    totalTrucks: 0,
    totalRoutes: 0,
    deliveryPointsByCategory: {},
    trucksByCapacity: {},
    routesByStatus: {}
  });

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        // Fetch delivery points
        const deliveryPointsResponse = await axios.get('/api/delivery-points');
        const deliveryPoints = deliveryPointsResponse.data;
        
        // Fetch trucks
        const trucksResponse = await axios.get('/api/trucks');
        const trucks = trucksResponse.data;
        
        // Fetch routes
        const routesResponse = await axios.get('/api/routes');
        const routes = routesResponse.data;
        
        // Process data for statistics
        const deliveryPointsByCategory = deliveryPoints.reduce((acc, point) => {
          acc[point.category] = (acc[point.category] || 0) + 1;
          return acc;
        }, {});
        
        const trucksByCapacity = trucks.reduce((acc, truck) => {
          acc[truck.capacity] = (acc[truck.capacity] || 0) + 1;
          return acc;
        }, {});
        
        const routesByStatus = routes.reduce((acc, route) => {
          acc[route.status] = (acc[route.status] || 0) + 1;
          return acc;
        }, {});
        
        setStats({
          totalDeliveryPoints: deliveryPoints.length,
          totalTrucks: trucks.length,
          totalRoutes: routes.length,
          deliveryPointsByCategory,
          trucksByCapacity,
          routesByStatus
        });
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, []);

  // Prepare chart data
  const deliveryPointsChartData = {
    labels: Object.keys(stats.deliveryPointsByCategory),
    datasets: [
      {
        label: 'Delivery Points by Category',
        data: Object.values(stats.deliveryPointsByCategory),
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

  const trucksChartData = {
    labels: Object.keys(stats.trucksByCapacity).map(cap => `${cap} Pallets`),
    datasets: [
      {
        label: 'Trucks by Capacity',
        data: Object.values(stats.trucksByCapacity),
        backgroundColor: 'rgba(255, 99, 132, 0.6)',
        borderColor: 'rgba(255, 99, 132, 1)',
        borderWidth: 1,
      },
    ],
  };

  const routesChartData = {
    labels: Object.keys(stats.routesByStatus),
    datasets: [
      {
        label: 'Routes by Status',
        data: Object.values(stats.routesByStatus),
        backgroundColor: [
          'rgba(255, 99, 132, 0.6)',
          'rgba(54, 162, 235, 0.6)',
          'rgba(75, 192, 192, 0.6)',
          'rgba(255, 206, 86, 0.6)',
        ],
        borderColor: [
          'rgba(255, 99, 132, 1)',
          'rgba(54, 162, 235, 1)',
          'rgba(75, 192, 192, 1)',
          'rgba(255, 206, 86, 1)',
        ],
        borderWidth: 1,
      },
    ],
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
        Dashboard
      </Typography>
      
      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Delivery Points
              </Typography>
              <Typography variant="h3">
                {stats.totalDeliveryPoints}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Trucks
              </Typography>
              <Typography variant="h3">
                {stats.totalTrucks}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                Routes
              </Typography>
              <Typography variant="h3">
                {stats.totalRoutes}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* Charts */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Delivery Points by Category
            </Typography>
            <Box sx={{ height: 300, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
              <Pie data={deliveryPointsChartData} options={{ maintainAspectRatio: false }} />
            </Box>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Trucks by Capacity
            </Typography>
            <Box sx={{ height: 300, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
              <Bar 
                data={trucksChartData} 
                options={{
                  maintainAspectRatio: false,
                  scales: {
                    y: {
                      beginAtZero: true,
                      ticks: {
                        precision: 0
                      }
                    }
                  }
                }} 
              />
            </Box>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Routes by Status
            </Typography>
            <Box sx={{ height: 300, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
              <Pie data={routesChartData} options={{ maintainAspectRatio: false }} />
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;