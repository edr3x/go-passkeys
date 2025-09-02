package handlers

import "github.com/edr3x/passkeys/internal/services"

type AuthHandlers struct {
	services *services.AuthSevice
}

func NewAuthHandlers(authSvc *services.AuthSevice) *AuthHandlers {
	return &AuthHandlers{
		services: authSvc,
	}
}
