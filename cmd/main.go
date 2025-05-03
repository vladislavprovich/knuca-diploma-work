package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/vladislavprovich/knuca-diploma-work/internal/database"
	"github.com/vladislavprovich/knuca-diploma-work/internal/models"
	"github.com/vladislavprovich/knuca-diploma-work/internal/repository"
	"github.com/vladislavprovich/knuca-diploma-work/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	// Skip React build check as requested by user
	log.Println("Skipping React build check, using existing frontend files")

	// Get MongoDB connection details from environment variables
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "logistics_db")

	// Connect to MongoDB
	mongoDB, err := database.NewMongoDB(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize repositories
	deliveryRepo := repository.NewDeliveryPointRepository(mongoDB.DB)
	truckRepo := repository.NewTruckRepository(mongoDB.DB)
	routeRepo := repository.NewRouteRepository(mongoDB.DB)

	// Initialize services
	routeOptimizerService := service.NewRouteOptimizerService(deliveryRepo, truckRepo, routeRepo)

	// Initialize data if needed
	ctx := context.Background()
	initializeData(ctx, deliveryRepo, truckRepo)

	// Initialize Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(corsMiddleware())

	// Serve static files for the frontend
	router.Static("/static", "./web/build/static")
	router.StaticFile("/", "./web/build/index.html")
	router.StaticFile("/favicon.ico", "./web/build/favicon.ico")

	// Handle React router paths - serve index.html for all client-side routes
	router.NoRoute(func(c *gin.Context) {
		// Check if the request is for an API endpoint
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// For all other routes, serve the React app's index.html
		c.File("./web/build/index.html")
	})

	// Initialize API handlers
	api := router.Group("/api")
	{
		// Delivery point routes
		api.GET("/delivery-points", handleGetAllDeliveryPoints(deliveryRepo))
		api.GET("/delivery-points/:id", handleGetDeliveryPointByID(deliveryRepo))
		api.PUT("/delivery-points/:id/pallets", handleUpdateDeliveryPointPallets(deliveryRepo))

		// Truck routes
		api.GET("/trucks", handleGetAllTrucks(truckRepo))
		api.GET("/trucks/:id", handleGetTruckByID(truckRepo))
		api.POST("/trucks", handleCreateTruck(truckRepo))
		api.DELETE("/trucks/:id", handleDeleteTruck(truckRepo))

		// Route routes
		api.GET("/routes", handleGetAllRoutes(routeRepo))
		api.GET("/routes/:id", handleGetRouteByID(routeRepo))
		api.POST("/routes/optimize", handleOptimizeRoutes(routeOptimizerService))
		api.GET("/routes/date/:date", handleGetRoutesByDate(routeRepo))
		api.PUT("/routes/:id", handleUpdateRoute(routeOptimizerService))
		api.DELETE("/routes", handleDeleteAllRoutes(routeOptimizerService))
		api.DELETE("/routes/:id", handleDeleteRoute(routeOptimizerService))
	}

	// Start the server
	port := getEnv("PORT", "6060") // Use environment variable with fallback
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server starting on port %s...", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	log.Printf("Received shutdown signal: %s", sign)
	defer func() {
		err = server.Shutdown(ctx)
	}()

	log.Print("Server shutdown complete")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// initializeData initializes the database with sample data if it's empty
func initializeData(ctx context.Context, deliveryRepo *repository.DeliveryPointRepository, truckRepo *repository.TruckRepository) {
	// Initialize delivery points with random distances
	if err := deliveryRepo.GenerateRandomDeliveryPoints(ctx); err != nil {
		log.Printf("Failed to initialize delivery points: %v", err)
	}

	// Initialize truck fleet
	if err := truckRepo.InitializeFleet(ctx); err != nil {
		log.Printf("Failed to initialize truck fleet: %v", err)
	}
}

// Handler functions
func handleGetAllDeliveryPoints(repo *repository.DeliveryPointRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		points, err := repo.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, points)
	}
}

func handleGetDeliveryPointByID(repo *repository.DeliveryPointRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pointID := 0
		if _, err := fmt.Sscanf(id, "%d", &pointID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		point, err := repo.GetByID(c.Request.Context(), pointID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, point)
	}
}

func handleUpdateDeliveryPointPallets(repo *repository.DeliveryPointRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pointID := 0
		if _, err := fmt.Sscanf(id, "%d", &pointID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var request struct {
			Pallets int `json:"pallets"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		point, err := repo.GetByID(c.Request.Context(), pointID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		point.Pallets = request.Pallets
		if err := repo.Update(c.Request.Context(), point); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, point)
	}
}

func handleGetAllTrucks(repo *repository.TruckRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		trucks, err := repo.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, trucks)
	}
}

func handleGetTruckByID(repo *repository.TruckRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		truckID := 0
		if _, err := fmt.Sscanf(id, "%d", &truckID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		truck, err := repo.GetByID(c.Request.Context(), truckID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, truck)
	}
}

func handleGetAllRoutes(repo *repository.RouteRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		routes, err := repo.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, routes)
	}
}

func handleGetRouteByID(repo *repository.RouteRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		routeID := 0
		if _, err := fmt.Sscanf(id, "%d", &routeID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		route, err := repo.GetByID(c.Request.Context(), routeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, route)
	}
}

func handleGetRoutesByDate(repo *repository.RouteRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}

		routes, err := repo.GetByDate(c.Request.Context(), date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, routes)
	}
}

// handleUpdateRoute handles updating an existing route
func handleUpdateRoute(service *service.RouteOptimizerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get route ID from URL parameter
		id := c.Param("id")
		routeID := 0
		if _, err := fmt.Sscanf(id, "%d", &routeID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		// Bind JSON request body to Route struct
		var route models.Route
		if err := c.ShouldBindJSON(&route); err != nil {
			log.Printf("Error binding JSON in route update: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		// Ensure the route ID in the URL matches the one in the request body
		if route.ID != routeID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Route ID in URL does not match ID in request body"})
			return
		}

		// Update the route using the service
		log.Printf("Updating route %d", routeID)
		if err := service.UpdateRoute(c.Request.Context(), &route); err != nil {
			log.Printf("Error updating route: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route: " + err.Error()})
			return
		}

		// Return the updated route
		c.JSON(http.StatusOK, route)
	}
}

func handleOptimizeRoutes(service *service.RouteOptimizerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			DeliveryDate string `json:"delivery_date"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("Error binding JSON in route optimization: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		if request.DeliveryDate == "" {
			log.Println("Error: Empty delivery date provided")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Delivery date is required"})
			return
		}

		date, err := time.Parse("2006-01-02", request.DeliveryDate)
		if err != nil {
			log.Printf("Error parsing date '%s': %v", request.DeliveryDate, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}

		log.Printf("Optimizing routes for date: %s", date.Format("2006-01-02"))
		routes, err := service.OptimizeRoutes(c.Request.Context(), date)
		if err != nil {
			log.Printf("Error optimizing routes: %v", err)

			// Check for specific error messages and provide more helpful responses
			if strings.Contains(err.Error(), "no suitable trucks available") ||
				strings.Contains(err.Error(), "no trucks available") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "No suitable trucks available for delivery. Please ensure trucks are available and properly configured.",
					"details": err.Error(),
				})
				return
			}

			if strings.Contains(err.Error(), "no delivery points with pallets") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "No delivery points with pallets found. Please add pallets to delivery points before optimizing routes.",
					"details": err.Error(),
				})
				return
			}

			// Generic error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to optimize routes",
				"details": err.Error(),
			})
			return
		}

		log.Printf("Successfully optimized %d routes, now saving", len(routes))
		// Save the optimized routes
		if err := service.SaveRoutes(c.Request.Context(), routes); err != nil {
			log.Printf("Error saving optimized routes: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save routes: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, routes)
	}
}

func handleDeleteAllRoutes(service *service.RouteOptimizerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := service.DeleteAllRoutes(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "All routes deleted successfully"})
	}
}

func handleCreateTruck(repo *repository.TruckRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var truck models.Truck
		if err := c.ShouldBindJSON(&truck); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate truck capacity
		if !models.IsValidCapacity(truck.Capacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid truck capacity"})
			return
		}

		// Set trailer information based on capacity
		truck.HasTrailer = truck.Capacity == 33
		if truck.HasTrailer && truck.TrailerInfo == "" {
			truck.TrailerInfo = fmt.Sprintf("Trailer-%d", truck.ID)
		}

		// Get next truck ID
		nextID, err := repo.GetNextTruckID(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		truck.ID = nextID

		// Set truck as available by default
		truck.Available = true

		// Create the truck
		if err := repo.Create(c.Request.Context(), &truck); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, truck)
	}
}

func handleDeleteTruck(repo *repository.TruckRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		truckID := 0
		if _, err := fmt.Sscanf(id, "%d", &truckID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		// Check if truck exists
		_, err := repo.GetByID(c.Request.Context(), truckID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Delete the truck
		if err := repo.Delete(c.Request.Context(), truckID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Truck deleted successfully"})
	}
}

// handleDeleteRoute handles the deletion of a specific route
func handleDeleteRoute(service *service.RouteOptimizerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		routeID := 0
		if _, err := fmt.Sscanf(id, "%d", &routeID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		if err := service.DeleteRoute(c.Request.Context(), routeID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Route deleted successfully"})
	}
}

// ensureReactBuild checks if the React build exists and builds it if needed
func ensureReactBuild() error {
	// Check if build directory exists
	buildDir := "./web/build"
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		log.Println("React build not found, building React app...")
		return buildReactApp()
	} else if err != nil {
		return fmt.Errorf("failed to check React build directory: %w", err)
	}

	// Check if index.html exists in build directory
	indexPath := filepath.Join(buildDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		log.Println("React build incomplete, rebuilding React app...")
		return buildReactApp()
	} else if err != nil {
		return fmt.Errorf("failed to check React build index.html: %w", err)
	}

	log.Println("React build found, using existing build")
	return nil
}

// buildReactApp builds the React app using npm
func buildReactApp() error {
	log.Println("Building React app...")

	// Change to web directory
	cmd := exec.Command("cmd", "/c", "cd web && npm run build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build React app: %w", err)
	}

	log.Println("React app built successfully")
	return nil
}
