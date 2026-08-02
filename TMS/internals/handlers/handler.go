package handlers

import "github.com/Jisnu-Dev/TMS/internals/services"

// Handler implements the gRPC TokenServiceServer interface.
// The generated interface from proto will be embedded here after running protoc.
type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}
