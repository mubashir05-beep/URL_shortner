package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mubashir05-beep/url_shortner/config"
	"github.com/mubashir05-beep/url_shortner/models"
	"github.com/mubashir05-beep/url_shortner/utils"
	"gorm.io/gorm/clause"
)

// Register user
func RegisterUser(c *fiber.Ctx) error {

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var missingFields []string
	if user.Name == "" {
		missingFields = append(missingFields, "Name")
	}
	if user.Email == "" {
		missingFields = append(missingFields, "Email")
	}
	if user.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	if len(missingFields) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("The following fields are required: %s", strings.Join(missingFields, ", ")),
		})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}
	user.Password = string(hashedPassword)
	err = config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error

	if err != nil {
		log.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if user.ID == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// LoginUser handles user login and JWT generation
func LoginUser(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate required fields
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	var user models.User

	// Fetch user from database
	err := config.DB.Raw("SELECT id, password FROM users WHERE email = ?", input.Email).Scan(&user).Error
	if err != nil {
		log.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Check if user exists
	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// Compare hashed password
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Println("JWT generation error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Return the token
	return c.JSON(fiber.Map{"token": token})
}
