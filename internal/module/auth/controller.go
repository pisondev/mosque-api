package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	GoogleLogin(c *fiber.Ctx) error
}

type controller struct {
	service Service
	log     *logrus.Logger
}

func NewController(service Service, log *logrus.Logger) Controller {
	return &controller{service: service, log: log}
}

func (ctrl *controller) GoogleLogin(c *fiber.Ctx) error {
	ctrl.log.Info("received google login request")

	var req LoginGoogleRequest
	if err := c.BodyParser(&req); err != nil {
		ctrl.log.Error("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request format"})
	}

	resp, err := ctrl.service.HandleGoogleLogin(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
