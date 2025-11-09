package checkinout

import (
	"github.com/gofiber/fiber/v2"
)

type CheckinoutRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterRoutes(app *fiber.App) {
	app.Post("/api/checkinout", Handler)
}

func Handler(c *fiber.Ctx) error {
	var req CheckinoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	username := req.Username
	password := req.Password
	if username == "" || password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "error": "username and password are required"})
	}

	err := Run(username, password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Check-in/out completed"})
}
