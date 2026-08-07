package handlers

import (
	"errors"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

var (
	ErrUnableToRegisterAdmin = errors.New("unable to register admin")
	ErrUnableToLoginAdmin    = errors.New("unable to login admin")

	ErrInvalidPayload = errors.New("invalid request payload")
	ErrInvalidIDParam = errors.New("invalid id parameter")
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
