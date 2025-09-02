package services

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/edr3x/passkeys/internal/db/sqlc"
)

type AuthSevice struct {
	webAuthn *webauthn.WebAuthn
	db       *pgxpool.Pool
	query    sqlc.Querier
	redis    *redis.Client
}

func getEnv(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}

func NewAuthService(db *pgxpool.Pool, redis *redis.Client) *AuthSevice {
	proto := getEnv("PROTOCOL", "http")
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", ":8080")
	origin := fmt.Sprintf("%s://%s%s", proto, host, port)

	wconfig := &webauthn.Config{
		RPDisplayName: "Passkey Webauthn",
		RPID:          host,
		RPOrigins:     []string{origin},
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		log.Fatal("Failed initializing:", err.Error())
		return nil
	}

	return &AuthSevice{
		webAuthn: webAuthn,
		db:       db,
		query:    sqlc.New(db),
		redis:    redis,
	}
}
