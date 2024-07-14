package controller

import (
	fiber "github.com/gofiber/fiber/v2"
	"time"
)

// Meta holds project metadata
type Meta struct {
	CommitHash string    `json:"commitHash"`
	BuildDate  string    `json:"buildDate"`
	SystemTime time.Time `json:"systemTime"`
}

// Test endpoint
func MetaHandler(c *fiber.Ctx) error {
	currentTimestamp := time.Now().UTC()
	data := Meta{
		BuildDate:  cfg.Meta.BuildDate,
		CommitHash: cfg.Meta.CommitHash,
		SystemTime: currentTimestamp,
	}

	return c.Status(200).JSON(data)
}
