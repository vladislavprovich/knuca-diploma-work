import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

// Layout components
import Layout from './components/Layout';

// Page components
import Dashboard from './pages/Dashboard';
import DeliveryPoints from './pages/DeliveryPoints';
import Trucks from './pages/Trucks';
import RoutesPage from './pages/Routes';
import RouteOptimizer from './pages/RouteOptimizer';
import Reports from './pages/Reports';
import Map from './pages/Map';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#f5f5f5',
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/delivery-points" element={<DeliveryPoints />} />
          <Route path="/trucks" element={<Trucks />} />
          <Route path="/routes" element={<RoutesPage />} />
          <Route path="/route-optimizer" element={<RouteOptimizer />} />
          <Route path="/reports" element={<Reports />} />
          <Route path="/map" element={<Map />} />
        </Routes>
      </Layout>
    </ThemeProvider>
  );
}

export default App;