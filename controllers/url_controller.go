package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mubashir05-beep/url_shortner/config"
	"github.com/mubashir05-beep/url_shortner/models"
	"github.com/mubashir05-beep/url_shortner/utils"
)

func AddURL(c *fiber.Ctx) error {
	var url models.URL

	if err := c.BodyParser(&url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if url.OriginalURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Original URL is required"})
	}
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	url.ShortCode = utils.GenerateShortCode()
	url.ClickCount = 0
	url.CreatedAt = time.Now()
	url.UserID = uint(userID)
	err := config.DB.Create(&url).Error
	if err != nil {
		log.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	return c.JSON(fiber.Map{"message": "URL shortened successfully", "short_url": url.ShortCode})

}

func ListURLs(c *fiber.Ctx) error {
	var urls []models.URL

	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	log.Println("User ID:", uint(userID))

	if err := config.DB.
		Preload("Analytics").Where("user_id = ?", uint(userID)).Find(&urls).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch URLs"})
	}

	return c.JSON(urls)
}

func DeleteURL(c *fiber.Ctx) error {
	shortCode := c.Params("short_code")
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var url models.URL
	if err := config.DB.Where("short_code = ? AND user_id = ?", shortCode, uint(userID)).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}
	if err := config.DB.Delete(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete URL"})
	}
	return c.JSON(fiber.Map{"message": "URL deleted successfully"})
}

func ViewAnalytics(c *fiber.Ctx) error {
	shortCode := c.Params("short_code")

	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}

	var analytics []models.Analytics
	if err := config.DB.Where("url_id = ?", url.ID).Find(&analytics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch analytics"})
	}

	return c.JSON(analytics)
}

func GetURLDetails(c *fiber.Ctx) error {

	shortCode := c.Params("short_code")

	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}

	var analytics []models.Analytics
	if err := config.DB.Where("url_id = ?", url.ID).Find(&analytics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch analytics"})
	}

	return c.JSON(analytics)
}

func GetURL(c *fiber.Ctx) error {
	shortCode := c.Params("short_code")

	// Find the URL
	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	}

	// Get client IP & user-agent details
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Parse user agent for browser & device info
	browser, device := utils.ParseUserAgent(userAgent)

	analytics := models.Analytics{
		URLID:     url.ID,
		IPAddress: ipAddress,
		Country:   utils.GetCountryFromIP(ipAddress), // Assuming you have a function to get country
		Device:    device,
		Browser:   browser,
		ClickedAt: time.Now(),
	}

	config.DB.Create(&analytics)

	// Increase click count
	url.ClickCount++
	config.DB.Save(&url)

	return c.Redirect(url.OriginalURL, fiber.StatusMovedPermanently)
}
