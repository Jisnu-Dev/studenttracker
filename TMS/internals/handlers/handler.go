package handlers

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/services"
)

type Handler struct {
	service services.ServiceInterface
	tokenpb.UnimplementedTokenServiceServer
}

func NewHandler(service services.ServiceInterface) *Handler {
	return &Handler{service: service}
}
