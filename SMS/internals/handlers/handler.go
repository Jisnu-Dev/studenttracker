package handlers

import "github.com/Jisnu-Dev/studenttracker/internals/services"

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}
