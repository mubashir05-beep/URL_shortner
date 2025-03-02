package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mubashir05-beep/url_shortner/controllers"
	"github.com/mubashir05-beep/url_shortner/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Public routes
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.LoginUser)

	// Protected API routes (Require authentication)
	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware()) // ✅ Apply auth middleware to protect routes

	api.Post("/urls", controllers.AddURL)                        // Create a short URL
	api.Get("/urls", controllers.ListURLs)                       // List all user's URLs
	api.Delete("/urls/:short_code", controllers.DeleteURL)       // Delete a specific URL
	api.Get("/analytics/:short_code", controllers.ViewAnalytics) // View analytics of a short URL
	api.Get("/details/:short_code", controllers.GetURLDetails)   // Get details of a specific short URL

	// Public redirection route
	app.Get("/:short_code", controllers.GetURL) // Redirect to original URL
}
