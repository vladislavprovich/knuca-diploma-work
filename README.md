# Logistics Management System

## Overview
This is a comprehensive logistics management system designed to optimize delivery routes for a fleet of trucks. The system allows for efficient planning and management of deliveries, with features for route optimization, truck management, and delivery point tracking.

## Features

### Route Optimization
- Automatically generates optimized delivery routes based on truck capacity and delivery point requirements
- Calculates total distance and cost for each route
- Supports different truck capacities and delivery point categories
- Implements category-based routing rules with flexibility for different truck types

### Truck Fleet Management
- Manages a fleet of trucks with different capacities (10, 15, 18, or 33 pallets)
- Tracks truck availability, driver information, and trailer details
- Calculates cost per kilometer based on truck capacity
- Supports trucks with and without trailers

### Delivery Point Tracking
- Categorizes delivery points by accessibility (Blue, Green, Yellow, Purple)
- Tracks pallet requirements for each delivery point
- Stores geographical information including address, latitude, and longitude
- Maintains distance data between delivery points for route calculation

### User Interface
- Modern web interface built with React
- Interactive route visualization
- Dashboard for monitoring delivery status
- Forms for managing trucks and delivery points

## Technology Stack

### Backend
- **Go (Golang)** - Main programming language
- **Gin** - Web framework
- **MongoDB** - Database
- **Docker** - Containerization

### Frontend
- **React** - UI library
- **JavaScript/CSS** - Frontend development

## Installation

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or higher
- Node.js and npm

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/your-username/knuca-diploma-work.git
   cd knuca-diploma-work
   ```

2. Set up environment variables:
   - Copy `.env.test` to `.env` and adjust the values as needed
   - Default configuration should work for local development

3. Start MongoDB using Docker:
   ```
   docker-compose up -d
   ```

4. Install frontend dependencies and build:
   ```
   cd web
   npm install
   npm run build
   cd ..
   ```

5. Run the application:
   ```
   powershell ./start.bat
   go run cmd/main.go
   ```
   
   Alternatively, use the provided batch file on Windows:
   ```
   start.bat
   ```

6. Access the application at http://localhost:6060

### Docker Deployment

To build and run the entire application using Docker:

```
# Build the Docker image
docker build -t logistics-app .

# Run the container
docker run -p 6060:6060 --network logistics-network logistics-app
```

## Usage

### Managing Trucks
- Add, edit, and delete trucks from the fleet
- Set truck capacity, driver information, and availability
- Track which trucks are currently assigned to routes

### Managing Delivery Points
- View all delivery points on the map
- Update pallet requirements for each delivery point
- View detailed information about each location

### Route Planning
- Generate optimized routes based on current delivery requirements
- View route details including distance, cost, and assigned truck
- Track route status (Planned, In Progress, Completed, Cancelled)

## Project Structure

```
├── cmd/                # Application entry point
│   └── main.go         # Main application file
├── internal/           # Internal packages
│   ├── database/       # Database connection and utilities
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic
├── web/                # Frontend React application
│   ├── public/         # Static assets
│   ├── src/            # React source code
│   └── build/          # Compiled frontend
├── docker-compose.yml  # Docker Compose configuration
├── Dockerfile          # Docker build configuration
└── .env                # Environment variables
```

## License

2025 Kyiv National University of Construction and Architecture