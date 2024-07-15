package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func BlackHole(c *fiber.Ctx) error {
	zap.L().Warn(
		fmt.Sprintf("Request caught by blackhole at path=%s", c.Path()),
	)

	return c.Status(404).SendString("Lol.")
}
