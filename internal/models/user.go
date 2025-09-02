package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/edr3x/passkeys/internal/db/sqlc"
)

type IPassKeyUser interface {
	webauthn.User

	AddCredential(*webauthn.Credential) error
	UpdateCredential(*webauthn.Credential) error
}

type PassKeyUser struct {
	id          uuid.UUID
	displayName string
	name        string
	email       string

	query sqlc.Querier
	redis *redis.Client
}

func GetPassKeyUser(
	id uuid.UUID,
	first_name,
	last_name,
	email string,
	query sqlc.Querier,
	redis *redis.Client,
) *PassKeyUser {
	return &PassKeyUser{
		id:          id,
		displayName: fmt.Sprintf("%s %s", first_name, last_name),
		name:        first_name,
		email:       email,
		query:       query,
		redis:       redis,
	}
}

func (u *PassKeyUser) WebAuthnID() []byte {
	return []byte(u.id.String())
}

func (u *PassKeyUser) WebAuthnName() string {
	return u.name
}

func (u *PassKeyUser) WebAuthnDisplayName() string {
	return u.displayName
}

func (u *PassKeyUser) WebAuthnCredentials() []webauthn.Credential {
	creds, err := u.query.GetUserPasskeyCredentials(context.Background(), u.id)
	if err != nil {
		return nil
	}

	webCreds := make([]webauthn.Credential, 0, len(creds))
	for _, c := range creds {
		// --- Transport ---
		transports := make([]protocol.AuthenticatorTransport, len(c.Transports))
		for i, t := range c.Transports {
			transports[i] = protocol.AuthenticatorTransport(t)
		}

		// --- Flags ---
		var flags webauthn.CredentialFlags
		if len(c.Flags) > 0 {
			if err := json.Unmarshal(c.Flags, &flags); err != nil {
				// return nil, fmt.Errorf("credential[%d] invalid flags: %w", idx, err)
				return nil
			}
		}

		// --- Attestation ---
		var attestation webauthn.CredentialAttestation
		if len(c.Attestation) > 0 {
			if err := json.Unmarshal(c.Attestation, &attestation); err != nil {
				// return nil, fmt.Errorf("credential[%d] invalid attestation: %w", idx, err)
				return nil
			}
		}

		// --- Build Credential ---
		cred := webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			Transport:       transports,
			AttestationType: c.AttestationType.String,
			Flags:           flags,
			Attestation:     attestation,
			Authenticator: webauthn.Authenticator{
				AAGUID:       c.Aaguid.UUID[:],
				SignCount:    uint32(c.SignCount),
				CloneWarning: c.CloneWarning.Bool,
				Attachment:   protocol.AuthenticatorAttachment(c.Attachment.String),
			},
		}

		webCreds = append(webCreds, cred)
	}
	return webCreds
}

func (o *PassKeyUser) AddCredential(credential *webauthn.Credential) error {
	transports := make([]string, len(credential.Transport))
	for i, t := range credential.Transport {
		transports[i] = string(t)
	}

	aaguid := func() uuid.NullUUID {
		id, err := uuid.Parse(string(credential.Authenticator.AAGUID))
		if err != nil {
			return uuid.NullUUID{Valid: false}
		}
		return uuid.NullUUID{Valid: true, UUID: id}
	}()

	attensation, _ := json.Marshal(credential.Attestation)

	flags, _ := json.Marshal(credential.Flags)

	arg := sqlc.AddCredentialParams{
		ID:              credential.ID,
		UserID:          o.id,
		PublicKey:       credential.PublicKey,
		SignCount:       int64(credential.Authenticator.SignCount),
		Transports:      transports,
		AttestationType: pgtype.Text{Valid: true, String: credential.AttestationType},
		Aaguid:          aaguid,
		Attestation:     attensation,
		Flags:           flags,
		CloneWarning:    pgtype.Bool{Valid: true, Bool: credential.Authenticator.CloneWarning},
		Attachment:      pgtype.Text{Valid: true, String: string(credential.Authenticator.Attachment)},
	}

	if err := o.query.AddCredential(context.Background(), arg); err != nil {
		return err
	}

	return nil
}

func (o *PassKeyUser) UpdateCredential(credential *webauthn.Credential) error {
	err := o.query.UpdateCredential(context.Background(), sqlc.UpdateCredentialParams{
		ID:        credential.ID,
		SignCount: int64(credential.Authenticator.SignCount),
	})
	if err != nil {
		return err
	}
	return nil
}
