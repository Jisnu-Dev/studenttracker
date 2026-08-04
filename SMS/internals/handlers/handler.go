package handlers

import (
	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

type Handler struct {
	service     services.ServiceInterface
	TokenClient grpcClient.TokenClientInterface
}

func NewHandler(service services.ServiceInterface, tokenClient grpcClient.TokenClientInterface) *Handler {
	return &Handler{
		service:     service,
		TokenClient: tokenClient,
	}
}
