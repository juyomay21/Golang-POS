package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Base struct {
	Trace string
}

func (b *Base) Auth(c *fiber.Ctx) error {
	if len(c.Get("Authorization")) > 0 {
		return nil
	} else {
		return errors.New("Auth Failed")
	}
}
