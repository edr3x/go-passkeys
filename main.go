package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"

	"github.com/edr3x/passkeys/internal/db/connection"
	"github.com/edr3x/passkeys/internal/handlers"
	"github.com/edr3x/passkeys/internal/redisclient"
	"github.com/edr3x/passkeys/internal/routes"
	"github.com/edr3x/passkeys/internal/services"
)

func main() {
	conn := connection.MustConnectPG(context.Background())
	defer conn.Close()
	redisClient := redisclient.RedisConnect()

	authSvc := services.NewAuthService(conn, redisClient)

	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"sid"},
	}))

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.Route("/api/auth", func(r chi.Router) {
		h := handlers.NewAuthHandlers(authSvc)
		routes.AuthRoutes(r, h)
	})

	mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "route not found", http.StatusNotFound)
	})

	http.ListenAndServe("0.0.0.0:8080", mux)
}
