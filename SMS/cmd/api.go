package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type AppConfig struct {
	ServerPort string
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

func mount(router *gin.Engine, handler *handlers.Handler) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "API is good."})
	})
	// Additional routes can be mounted here using the handler
}

func run() error {
	config := loadConfig()

	db, err := connectToDB(&config.DB)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	service := services.NewService(db)
	handler := handlers.NewHandler(service)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	mount(router, handler)

	log.Printf("Server starting on port %s", config.ServerPort)
	if err := router.Run(":" + config.ServerPort); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

func loadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found or error reading it, relying on environment variables")
	}

	return &AppConfig{
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
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

func connectToDB(cfg *DatabaseConfig) (*sql.DB, error) {
	dsn := getDSN(cfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to DB successfully")
	return db, nil
}
