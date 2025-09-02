package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"

	"github.com/edr3x/passkeys/internal/db/sqlc"
	"github.com/edr3x/passkeys/internal/models"
)

func (s *AuthSevice) InitializePasskey(ctx context.Context, email string) (any, error) {
	user, err := s.getPasskeyUser(ctx, email)
	if err != nil {
		return nil, err
	}

	// Resident key requirement
	opts := func(credOpts *protocol.PublicKeyCredentialCreationOptions) {
		credOpts.AuthenticatorSelection = protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationPreferred,
		}
	}

	options, session, err := s.webAuthn.BeginRegistration(user, opts)
	if err != nil {
		return nil, err
	}

	sid, _ := uuid.NewV7()

	sessionBytes, err := json.Marshal(*session)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(sid.String(), sessionBytes, 10*time.Minute).Err(); err != nil {
		log.Println("Setting error:", err.Error())
	}

	return map[string]any{
		"sid":  sid,
		"opts": options,
	}, nil
}

func (s *AuthSevice) FinalizePassKey(r *http.Request) (any, error) {
	sessionId := r.Header.Get("sid")
	ctx := r.Context()
	sessionStr, err := s.redis.Get(sessionId).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var session webauthn.SessionData
	if err := json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return nil, err
	}

	user, err := s.getPasskeyUser(ctx, string(session.UserID))
	if err != nil {
		return nil, err
	}

	creds, err := s.webAuthn.FinishRegistration(user, session, r)
	if err != nil {
		return nil, err
	}

	user.AddCredential(creds)

	s.redis.Del(sessionId)

	return "Success", nil
}

func (s *AuthSevice) BeginLogin(ctx context.Context) (any, error) {
	options, session, err := s.webAuthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, err
	}

	sid, _ := uuid.NewV7()

	sessionBytes, err := json.Marshal(*session)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(sid.String(), sessionBytes, 2*time.Minute).Err(); err != nil {
		log.Println("Setting error:", err.Error())
	}

	return map[string]any{
		"sid":  sid,
		"opts": options,
	}, nil
}

func (s *AuthSevice) FinishLogin(r *http.Request) (any, error) {
	sessionId := r.Header.Get("sid")

	ctx := r.Context()
	sessionStr, err := s.redis.Get(sessionId).Result()
	if err != nil {
		return nil, fmt.Errorf("not found")
	}

	var session webauthn.SessionData
	if err := json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return nil, err
	}

	var user *models.PassKeyUser
	discoverableUserHandler := func(_, userid []byte) (u webauthn.User, err error) {
		parsedUser, err := s.getPasskeyUser(ctx, string(userid))
		if err != nil {
			return nil, err
		}
		user = parsedUser
		return user, err
	}

	creds, err := s.webAuthn.FinishDiscoverableLogin(discoverableUserHandler, session, r)
	if err != nil {
		return nil, err
	}

	user.UpdateCredential(creds)

	s.redis.Del(sessionId)

	// return access token from here for furtner
	return "this is token string", nil
}

func (s *AuthSevice) getPasskeyUser(ctx context.Context, emailOrId string) (*models.PassKeyUser, error) {
	id, _ := uuid.Parse(emailOrId)
	arg := sqlc.GetUserByEmailOrIdParams{
		Email: emailOrId,
		ID:    id,
	}
	pUser, err := s.query.GetUserByEmailOrId(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("query.GetUserMinDetail: %w", err)
	}
	u := models.GetPassKeyUser(
		pUser.ID,
		pUser.FirstName,
		pUser.LastName,
		pUser.Email,
		s.query,
		s.redis,
	)
	return u, nil
}
