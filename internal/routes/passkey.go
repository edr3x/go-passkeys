package routes

import (
	"github.com/edr3x/passkeys/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func PassKeyHandlers(r chi.Router, h *handlers.AuthHandlers) {
	r.Post("/initialize-start", h.InitializePasskey)
	r.Post("/initialize-finish", h.FinalizePasskey)
	r.Post("/login-start", h.BeginLogin)
	r.Post("/login-finish", h.FinishLogin)
}
