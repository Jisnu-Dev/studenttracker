package main

import (
	"fmt"
	"log"
	"net"
	"os"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/handlers"
	"github.com/Jisnu-Dev/TMS/internals/services"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type AppConfig struct {
	GRPCPort  string
	JWTSecret string
}

func run() error {
	config := loadConfig()

	service := services.NewService(config.JWTSecret)
	handler := handlers.NewHandler(service)

	grpcServer := grpc.NewServer()

	tokenpb.RegisterTokenServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	log.Printf("TMS gRPC server starting on port %s", config.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("gRPC server failed: %w", err)
	}

	return nil
}

func loadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found, relying on environment variables")
	}

	return &AppConfig{
		GRPCPort:  getEnvOrDefault("GRPC_PORT", "50051"),
		JWTSecret: getEnvOrDefault("JWT_SECRET", "d8f34a91c7e26b5f8a1d9c4e72bf06a5e1c398fd4b7a2e6c90f15d8ab34c7e29"),
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
