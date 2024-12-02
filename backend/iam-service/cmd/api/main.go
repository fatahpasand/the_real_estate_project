// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"iam-service/internal/adapters/handlers"
	mysqlrepo "iam-service/internal/adapters/repositories/mysql"
	redisrepo "iam-service/internal/adapters/repositories/redis"
	"iam-service/internal/adapters/services"
	"iam-service/internal/core/usecases"
	"iam-service/internal/middleware"
)

func main() {
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Initialize Fiber app with config
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		DisableStartupMessage: true,
	})

	// Add static file handling
	app.Static("/", "./static")

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize dependencies
	db, err := initializeMySQL()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	redisClient, err := initializeRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo := mysqlrepo.NewMySQLUserRepository(db)
	cacheRepo := redisrepo.NewRedisRepository(redisClient)
	auditRepo := mysqlrepo.NewMySQLAuditRepository(db)

	// Initialize services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(cacheRepo)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, cacheRepo, auditRepo, emailService, authService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)

	// Handle favicon specifically
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		// Return 204 No Content if no favicon exists
		return c.SendStatus(fiber.StatusNoContent)
	})

	// Get allowed origins from env or use default
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000" // Default for development
	}

	// Configure CORS with secure settings
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Add global middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(
		redisClient,
		60,          // 60 requests
		time.Minute, // per minute
	)

	// Add rate limiter middleware
	app.Use(rateLimiter.Middleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API routes
	api := app.Group("/api/v1")

	// Public routes with validation
	api.Post("/register", middleware.ValidateRegister, userHandler.Register)
	api.Post("/login", middleware.ValidateLogin, userHandler.Login)
	api.Get("/verify/:token", userHandler.VerifyEmail)

	// Documentation routes
	api.Get("/docs/openapi.yaml", middleware.ServeOpenAPISpec)
	api.Get("/docs", middleware.ServeRedoc)

	// Protected routes
	protected := api.Group("/", middleware.JWTAuth())
	protected.Get("/profile", userHandler.GetProfile)
	protected.Put("/profile", userHandler.UpdateProfile)

	// Add root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "IAM Service",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Documentation route
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "IAM Service API",
			"version": "v1",
			"endpoints": []string{
				"POST /register - Register new user",
				"POST /login - User login",
				"GET /verify/:token - Verify email",
				"GET /oauth/google - Google OAuth login",
				"GET /oauth/google/callback - Google OAuth callback",
				"GET /profile - Get user profile (Protected)",
				"PUT /profile - Update user profile (Protected)",
			},
		})
	})

	// Setup routes
	api = app.Group("/api/v1")

	// Public routes
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Get("/verify/:token", userHandler.VerifyEmail)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initializeRedis() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		return nil, fmt.Errorf("missing required Redis configuration")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}

func initializeMySQL() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Validate required fields
	if dbUser == "" || dbPassword == "" || dbName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	log.Printf("Connecting to MySQL at %s:%s...", dbHost, dbPort)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

func loadEnv() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Set required environment variables with defaults
	requiredEnvs := map[string]string{
		"DB_USER":     "root",
		"DB_PASSWORD": "root",
		"DB_HOST":     "localhost",
		"DB_PORT":     "3306",
		"DB_NAME":     "iam_db",
	}

	for key, defaultValue := range requiredEnvs {
		if os.Getenv(key) == "" {
			os.Setenv(key, defaultValue)
		}
	}

	return nil
}
