package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	connection *pgx.Conn
	redis      *redis.Client
}

func NewHealthHandler(conn *pgx.Conn, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		connection: conn,
		redis:      redis,
	}
}

func (h *HealthHandler) Handle(c *fiber.Ctx) error {
	if err := h.connection.Ping(c.Context()); err != nil {
		log.Error("connection.Ping err: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Database is not OK")
	}

	if err := h.redis.Ping(c.Context()).Err(); err != nil {
		log.Error("redis.Ping err: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Redis is not OK")
	}

	return c.Status(fiber.StatusOK).SendString("OK")
}
