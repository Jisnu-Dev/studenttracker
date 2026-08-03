package handlers_test

import (
	"github.com/Jisnu-Dev/studenttracker/internals/clients"
	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

// newHandler creates a Handler wired to the given mock dependencies for use in tests.
// Pass nil for tokenClient in handlers that do not interact with it directly.
func newHandler(svc *mocks.MockService, tc clients.TokenClientInterface) *handlers.Handler {
	return handlers.NewHandler(svc, tc)
}
