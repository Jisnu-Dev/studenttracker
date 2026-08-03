package handlers_test

import (
	"github.com/Jisnu-Dev/studenttracker/internals/clients"
	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func newHandler(svc *mocks.MockService, tc clients.TokenClientInterface) *handlers.Handler {
	return handlers.NewHandler(svc, tc)
}
