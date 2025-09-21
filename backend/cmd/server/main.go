package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"news-to-text/internal/config"
	"news-to-text/internal/database"
	"news-to-text/internal/cache"
	"news-to-text/internal/handlers"
	"news-to-text/internal/middleware"
	"news-to-text/internal/services"
	"news-to-text/internal/repositories"
	"news-to-text/pkg/auth"
	"news-to-text/pkg/logger"
	_ "news-to-text/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title News to Text API
// @version 1.0
// @description A news alert system that sends text notifications
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis cache
	redisClient, err := cache.Initialize(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	alertRepo := repositories.NewAlertRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, redisClient, cfg.JWTSecret)
	alertService := services.NewAlertService(alertRepo, redisClient)
	newsService := services.NewNewsService(cfg.NewsAPIKey)
	notificationService := services.NewNotificationService(cfg.SMSAPIKey)

	// Initialize JWT manager for middleware
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	alertHandler := handlers.NewAlertHandler(alertService, authService)

	// Initialize background services
	backgroundService := services.NewBackgroundService(alertService, newsService, notificationService)
	go backgroundService.Start()

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// Alert routes (protected)
		alerts := v1.Group("/alerts")
		alerts.Use(middleware.AuthMiddleware(jwtManager))
		{
			alerts.GET("", alertHandler.GetAlerts)
			alerts.POST("", alertHandler.CreateAlert)
			alerts.PUT("/:id", alertHandler.UpdateAlert)
			alerts.DELETE("/:id", alertHandler.DeleteAlert)
			alerts.GET("/history", alertHandler.GetAlertHistory)
			alerts.POST("/test", alertHandler.TestAlert)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Stop background services
	backgroundService.Stop()

	// The context is used to inform the server it has 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}