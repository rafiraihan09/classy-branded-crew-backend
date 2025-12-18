package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors any) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
