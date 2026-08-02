package handlers

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/services"
)

type Handler struct {
	service *services.Service
	tokenpb.UnimplementedTokenServiceServer
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}
