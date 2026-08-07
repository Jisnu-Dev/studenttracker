package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/middlewares"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type AppConfig struct {
	ServerPort string
	TMSUrl     string
	DB         DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func mount(router *gin.Engine, handler *handlers.Handler, tokenClient *grpcClient.TokenClient) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "API is good."})
	})
	// admin routes
	router.POST("/register", handler.RegisterAdminHandler)
	router.POST("/login", handler.LoginAdminHandler)

	// student routes
	studentRoutes := router.Group("/students")
	studentRoutes.Use(middlewares.AuthMiddleware(tokenClient))

	studentRoutes.POST("/", handler.CreateStudentHandler)
	studentRoutes.GET("/", handler.GetAllStudentsHandler)
	studentRoutes.GET("/:id", handler.GetStudentByIDHandler)
	studentRoutes.PUT("/:id", handler.UpdateStudentHandler)
	studentRoutes.PATCH("/:id", handler.PatchStudentHandler)
	studentRoutes.DELETE("/:id", handler.DeleteStudentHandler)
}

func run() {
	config := loadConfig()

	db := connectToDB(&config.DB)
	defer db.Close()

	//connecting to TMS gRPC service
	tokenClient, err := grpcClient.NewTokenClient(config.TMSUrl)
	if err != nil {
		slog.Error("Failed to create token client", "error", err, "function", "run", "location", "cmd/api.go")
		os.Exit(1)
	}

	service := services.NewService(db)
	handler := handlers.NewHandler(service, tokenClient)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	mount(router, handler, tokenClient)

	log.Printf("Server starting on port %s", config.ServerPort)
	if err := router.Run(":" + config.ServerPort); err != nil {
		slog.Error("Server failed to start", "error", err, "function", "run", "location", "cmd/api.go")
		os.Exit(1)
	}
}

func loadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found or error reading it, relying on environment variables")
	}

	return &AppConfig{
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
		TMSUrl:     getEnvOrDefault("TMS_URL", "localhost:50051"),
		DB: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", "postgres"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		if value != "" {
			return value
		}
	}
	return fallback
}

func getDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func connectToDB(cfg *DatabaseConfig) *sql.DB {
	dsn := getDSN(cfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err, "function", "connectToDB", "location", "cmd/api.go", "host", cfg.Host)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err, "function", "connectToDB", "location", "cmd/api.go")
		os.Exit(1)
	}

	log.Println("Connected to DB successfully")
	return db
}
