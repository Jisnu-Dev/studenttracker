package handlers

import (
	"github.com/Jisnu-Dev/studenttracker/internals/clients"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

type Handler struct {
	service     *services.Service
	TokenClient *clients.TokenClient
}

func NewHandler(service *services.Service, tokenClient *clients.TokenClient) *Handler {
	return &Handler{
		service:     service,
		TokenClient: tokenClient,
	}
}
