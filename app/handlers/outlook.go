package handlers

import (
	"encoding/xml"
	"github.com/gofiber/fiber/v2"
	"network-it/go-mail-autoconf/app/template"
)

func OutlookHandler(c *fiber.Ctx) error {
	type AutoDiscover struct {
		XMLName xml.Name `xml:"Autodiscover"`
		Text    string   `xml:",chardata"`
		Xmlns   string   `xml:"xmlns,attr"`
		Request struct {
			Text                     string `xml:",chardata"`
			EMailAddress             string `xml:"EMailAddress"`
			AcceptableResponseSchema string `xml:"AcceptableResponseSchema"`
		} `xml:"Request"`
	}

	//type AutoDiscover struct {
	//	AutoDiscover struct {
	//		Request struct {
	//			EMailAddress string `xml:"EMailAddress"`
	//		} `xml:"Request"`
	//	} `xml:"Autodiscover"`
	//}
	var reqData AutoDiscover
	err := c.BodyParser(&reqData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if reqData.Request.EMailAddress == "" {
		return c.Status(fiber.StatusBadRequest).SendString("EMailAddress parameter is required")
	}

	data := template.TemplateData{Address: reqData.Request.EMailAddress, Username: reqData.Request.EMailAddress, DisplayName: reqData.Request.EMailAddress}
	resData, err := template.Execute("outlook", data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationXMLCharsetUTF8)
	return c.SendString(resData)
}
