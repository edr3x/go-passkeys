package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/edr3x/passkeys/internal/utils"
)

type PasskeyParam struct {
	Email     string `json:"email"`
	SessionId string `json:"sid"`
}

func (h *AuthHandlers) InitializePasskey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body PasskeyParam
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	resp, err := h.services.InitializePasskey(ctx, body.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, resp)
}

func (h *AuthHandlers) FinalizePasskey(w http.ResponseWriter, r *http.Request) {
	resp, err := h.services.FinalizePassKey(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, map[string]any{
		"successs": true,
		"message":  resp,
	})
}

func (h *AuthHandlers) BeginLogin(w http.ResponseWriter, r *http.Request) {
	resp, err := h.services.BeginLogin(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, resp)
}

func (h *AuthHandlers) FinishLogin(w http.ResponseWriter, r *http.Request) {
	resp, err := h.services.FinishLogin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, resp)
}
