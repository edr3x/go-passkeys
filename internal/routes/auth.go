package routes

import (
	"github.com/edr3x/passkeys/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, handlers *handlers.AuthHandlers) {
	r.Route("/passkey", func(rt chi.Router) {
		PassKeyHandlers(rt, handlers)
	})
}
