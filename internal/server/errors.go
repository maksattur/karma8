package server

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	origin   error
	response *fiber.Error
}

func (e Error) Error() string {
	return fmt.Sprintf("http status %d: %s", e.response.Code, e.origin.Error())
}

func NewErrorsLoggerMiddleware(next fiber.ErrorHandler, logger logr.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		logger.Error(err, "request failed", "path", string(ctx.Request().URI().Path()))

		var wrapper *Error
		if errors.As(err, &wrapper) {
			err = wrapper.response
		}

		return next(ctx, err)
	}
}

func NewError(origin error, response *fiber.Error) *Error {
	return &Error{origin: origin, response: response}
}
