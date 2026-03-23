package response

import "github.com/gofiber/fiber/v2"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}, meta interface{}) error {
	if meta == nil {
		meta = fiber.Map{}
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

func Validation(c *fiber.Ctx, message string, errors []FieldError) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"errors":  errors,
	})
}
