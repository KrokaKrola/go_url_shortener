package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	App        *fiber.App
	connection *pgx.Conn
	redis      *redis.Client
}

func NewServer(connection *pgx.Conn, redis *redis.Client) *Server {
	return &Server{
		App:        fiber.New(),
		connection: connection,
		redis:      redis,
	}
}

func (s *Server) Start() error {
	InitMiddlewares(s.App)
	NewRoutes().Setup(s.App, s.connection, s.redis)

	log.Println("Server is running on port 3000")

	return s.App.Listen(":3000")
}
