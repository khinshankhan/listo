package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/khinshankhan/listo/internal/service/log"
)

func BlackHole(c *fiber.Ctx) error {
	logger := log.NewLogger()
	logger.Warn(
		fmt.Sprintf("Request caught by blackhole at path=%s", c.Path()),
	)

	return c.Status(404).SendString("Lol.")
}
