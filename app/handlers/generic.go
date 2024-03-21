package handlers

import (
	"github.com/gofiber/fiber/v2"
	"network-it/go-mail-autoconf/app/template"
	"strings"
)

// GenericConfigHandler is a handler for the generic configuration used for Thunderbird and iOS
func GenericConfigHandler(c *fiber.Ctx, templateKey string, paramName string) error {
	address := c.Query(paramName, "")
	if address == "" {
		return c.Status(fiber.StatusBadRequest).SendString(paramName + " parameter is required")
	}
	if !strings.Contains(address, "@") {
		return c.Status(fiber.StatusBadRequest).SendString(paramName + " parameter is not a valid email address")
	}
	data := template.NewData(address)
	resData, err := template.Execute(templateKey, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXMLCharsetUTF8)
	return c.SendString(resData)
}
