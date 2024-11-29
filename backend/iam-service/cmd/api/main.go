// cmd/api/main.go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "gorm.io/gorm"
    "github.com/go-redis/redis/v8"
)

func main() {
    app := fiber.New()
    app.Use(logger.New())

    // Initialize MySQL
    db, err := database.NewMySQLConnection(&database.Config{
        Host:     "localhost",
        Port:     3306,
        User:     "root",
        Password: "your-password",
        DBName:   "amlak",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:5891",
    })

    // Initialize repositories
    userRepo := repository.NewMySQLUserRepository(db)
    
    // Initialize services
    userUseCase := usecases.NewUserUseCase(userRepo, redisClient)
    
    // Setup routes
    handler := handlers.NewUserHandler(userUseCase)
    
    api := app.Group("/api/v1")
    api.Post("/register", handler.Register)
    api.Post("/login", handler.Login)
    api.Get("/verify/:token", handler.VerifyEmail)
    api.Get("/oauth/google", handler.GoogleOAuth)
    
    app.Listen(":3000")
}