package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type JsonRequestValidator struct {
}

type JsonResponseSender struct {
}

type FormRequestValidator struct {
}

func (j *JsonRequestValidator) ValidateJsonBody(c *fiber.Ctx, dto interface{}) error {

	if err := c.BodyParser(&dto); err != nil {
		return err

	}
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return err

	}

	return nil

}

func (j *JsonResponseSender) JSONResponse(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status": statusCode,
		"data":   data,
	})
}
