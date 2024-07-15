package controller

import (
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/khinshankhan/listo/internal/config"
	"go.uber.org/zap"
)

const (
	// TODO: move this to config?
	// Port is the port that the server should listen to
	Port = 8080
)

var (
	cfg *config.Config
)

type FiberFunc = func(c *fiber.Ctx) error

// CreateRouter creates the router for the http server
func CreateRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		ProxyHeader: "Cf-Connecting-Ip",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				switch err.Error() {
				case "Method Not Allowed":
					return c.Status(404).SendString("not found")
				default:
					fmt.Println("hi")
					fmt.Println(err.Error())
				}
			}

			zap.L().Error(
				"Handler didn't work",
				zap.String("context", "ErrorHandler"),
				zap.Error(err),
			)
			return fiber.DefaultErrorHandler(c, err)
		},
	})

	// Prevent API calls from being handled by the catchall
	app.Get("/api/*", BlackHole)
	app.Get("/api/v1/meta", MetaHandler)

	return app
}

// Handle creates the router and the server before starting the server
func Handle(loadedCfg *config.Config) {
	cfg = loadedCfg

	app := CreateRouter()
	// TODO: use env variables + only do this for local dev server
	// TODO: use prod url when in prod, maybe get url from config too
	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:3000"}))
	zap.L().DPanic(
		"Something went terribly wrong",
		zap.String("context", "Handle"),
		zap.Error(app.Listen(fmt.Sprintf(":%d", Port))),
	)
}
