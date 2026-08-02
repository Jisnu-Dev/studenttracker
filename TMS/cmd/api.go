package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Jisnu-Dev/TMS/internals/handlers"
	"github.com/Jisnu-Dev/TMS/internals/services"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type AppConfig struct {
	GRPCPort string
}

func run() error {
	config := loadConfig()

	service := services.NewService()
	handler := handlers.NewHandler(service)

	// until we register the handler below
	_ = handler

	grpcServer := grpc.NewServer()

	// TODO: Register your generated service here after running protoc:
	// tokenpb.RegisterTokenServiceServer(grpcServer, handler)

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
		GRPCPort: getEnvOrDefault("GRPC_PORT", "50051"),
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
