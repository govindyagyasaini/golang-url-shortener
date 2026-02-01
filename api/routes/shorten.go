package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/govindyagyasaini/golang-url-shortener/database"
	"github.com/govindyagyasaini/golang-url-shortener/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// -------- RATE LIMITING (DB 1) --------
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute)
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "rate limit exceeded",
				"rate_limit_reset": ttl.Minutes(),
			})
		}
	}

	// -------- URL VALIDATION --------
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid domain",
		})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	// -------- SHORT ID --------
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// -------- STORE URL (DB 0) --------
	r := database.CreateClient(0)
	defer r.Close()

	if existing, _ := r.Get(database.Ctx, id).Result(); existing != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "short URL already exists",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(
		database.Ctx,
		id,
		body.URL,
		body.Expiry*time.Hour,
	).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot save URL",
		})
	}

	// -------- RESPONSE --------
	r2.Decr(database.Ctx, c.IP())
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  atoi(val),
		XRateLimitReset: ttl,
	}

	return c.JSON(resp)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
