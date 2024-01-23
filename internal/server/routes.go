package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/krokakrola/url_shortener/internal/handlers"
	"github.com/redis/go-redis/v9"
)

type Routes struct{}

func NewRoutes() *Routes {
	return &Routes{}
}

func (r *Routes) Setup(app *fiber.App, connection *pgx.Conn, redis *redis.Client) {
	apiV1 := app.Group("/api/v1")

	validate := NewValidator().validate

	healthHandler := handlers.NewHealthHandler(connection, redis)
	shortenHandler := handlers.NewShortenHandler(connection, validate, redis)

	apiV1.Get("/health", healthHandler.Handle)
	apiV1.Post("/shorten", shortenHandler.HandleCreateShortURL)
	apiV1.Get("/:shortUrl", shortenHandler.HandleRedirect)
	apiV1.Get("/:shortUrl/stats", shortenHandler.HandleGetVisitsCount)
}
