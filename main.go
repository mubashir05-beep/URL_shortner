package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mubashir05-beep/url_shortner/config"
	"github.com/mubashir05-beep/url_shortner/routes"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":3000")
}
