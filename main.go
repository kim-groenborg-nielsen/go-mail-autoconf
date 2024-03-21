package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"network-it/go-mail-autoconf/app/config"
	"network-it/go-mail-autoconf/app/handlers"
	"network-it/go-mail-autoconf/app/template"
	"os"
)

var version string = ""
var commit string = ""
var date string = ""

func main() {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		log.Printf("Version: %s\nCommit: %s\nDate: %s", version, commit, date)
		os.Exit(0)
	}

	cfg := config.NewConfig()
	template.New()
	// Code
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${ips} | ${method} | ${path} | ${error}\n",
	}))

	// Thunderbird autoconfig
	app.Get("/mail/config-v1.1.xml", func(c *fiber.Ctx) error {
		return handlers.GenericConfigHandler(c, "thunderbird", "emailaddress")
	})

	// Outlook autodiscover
	app.Post("/autodiscover/autodiscover.xml", handlers.OutlookHandler)
	app.Post("/Autodiscover/Autodiscover.xml", handlers.OutlookHandler)
	app.Get("/autodiscover/autodiscover.xml", func(ctx *fiber.Ctx) error {
		return handlers.GenericConfigHandler(ctx, "outlook", "email")
	})
	app.Get("/autodiscover/autodiscover.json/v1.0/:email?", func(ctx *fiber.Ctx) error {
		url := fmt.Sprintf("%s://%s/autodiscover/autodiscover.xml", ctx.Get("x-forwarded-proto"), ctx.Get("host"))
		return ctx.JSON(fiber.Map{"Protocol": "Autodiscoverv1.0", "Url": url})
	})

	// iOS mobileconfig
	app.Get("/email.mobileconfig", func(c *fiber.Ctx) error {
		return handlers.GenericConfigHandler(c, "ios", "email")
	})

	if err := app.Listen(cfg.ServerUrl); err != nil {
		log.Fatal(err)
	}
}
