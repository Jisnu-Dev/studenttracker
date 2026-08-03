package handlers

import (
	"github.com/Jisnu-Dev/studenttracker/internals/clients"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

type Handler struct {
	service     services.ServiceInterface
	TokenClient clients.TokenClientInterface
}

func NewHandler(service services.ServiceInterface, tokenClient clients.TokenClientInterface) *Handler {
	return &Handler{
		service:     service,
		TokenClient: tokenClient,
	}
}
