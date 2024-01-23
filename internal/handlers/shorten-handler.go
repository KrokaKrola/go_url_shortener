package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/krokakrola/url_shortener/internal/store"
	"github.com/krokakrola/url_shortener/internal/utils"
	"github.com/redis/go-redis/v9"
)

type ShortenHandler struct {
	connection *pgx.Conn
	validate   *validator.Validate
	redis      *redis.Client
}

func NewShortenHandler(connection *pgx.Conn, validate *validator.Validate, redis *redis.Client) *ShortenHandler {
	return &ShortenHandler{
		connection: connection,
		validate:   validate,
		redis:      redis,
	}
}

func (h *ShortenHandler) HandleCreateShortURL(ctx *fiber.Ctx) error {
	req := store.NewCreateLinkRequest(ctx.Query("url"))

	if err := h.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": utils.ConvertValidationErrorsToObject(err),
		})
	}

	link := store.NewLink(req.URL, "")

	randomHash, err := link.CreateOne(ctx.Context(), h.connection)

	if err != nil {
		log.Error("link.CreateOne err: ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	response := store.NewCreateLinkResponse(randomHash)

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *ShortenHandler) HandleRedirect(ctx *fiber.Ctx) error {
	shortUrl := ctx.Params("shortUrl")

	if shortUrl == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short url is required",
		})
	}

	rValue, err := h.redis.Get(ctx.Context(), fmt.Sprintf("url:%s", shortUrl)).Result()

	var rResult string

	if err == redis.Nil {
		log.Info("Key does not exist")
	} else if err != nil {
		log.Error("Error getting key ", err)
	} else {
		rResult = rValue
	}

	var link *store.Link

	if rResult != "" {
		if err := json.Unmarshal([]byte(rResult), &link); err != nil {
			log.Error("Error unmarshaling link", err)

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Something went wrong",
			})
		}

		log.Info("Link got from Redis")
	} else {
		log.Info("Link got from DB")
		link, err = link.FindOneByHash(ctx.Context(), h.connection, shortUrl)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "link not found",
			})
		}

		log.Error("Error getting link", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	headers := ctx.GetReqHeaders()
	ip := ctx.IP()

	go func(ctx context.Context, headers map[string][]string, ip string) {
		var userAgent string
		if ua, ok := headers["User-Agent"]; ok && len(ua) > 0 {
			userAgent = ua[0]
		}

		var referer string
		if ref, ok := headers["Referer"]; ok && len(ref) > 0 {
			referer = ref[0]
		}

		visit := store.NewVisit(int(link.ID), userAgent, ip, referer)

		if err := visit.CreateOne(ctx, h.connection); err != nil {
			log.Info("Error writing visit", err)
		} else {
			log.Info("Visit written")
		}
	}(ctx.Context(), headers, ip)

	json, err := json.Marshal(link)

	if err != nil {
		log.Error("Error marshaling link", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	if h.redis.Set(ctx.Context(), fmt.Sprintf("url:%s", shortUrl), json, time.Second*600).Err(); err != nil {
		log.Error("Error setting key", err)
	}

	return ctx.Redirect(link.URL)
}

func (h *ShortenHandler) HandleGetVisitsCount(ctx *fiber.Ctx) error {
	shortUrl := ctx.Params("shortUrl")

	if shortUrl == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short url is required",
		})
	}

	var link *store.Link

	link, err := link.FindOneByHash(ctx.Context(), h.connection, shortUrl)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "link not found",
			})
		}

		log.Error("Error getting link", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	var visit *store.Visit = store.NewVisit(int(link.ID), "", "", "")
	visits, err := visit.GetVisitsCount(ctx.Context(), h.connection)

	if err != nil {
		log.Error("Error getting visits", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"visits": visits,
	})
}
