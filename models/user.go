package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"` // Hide password from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relation (One User -> Many URLs)
	Urls []URL `gorm:"foreignKey:UserID" json:"-"` // Hide URLs from user JSON response
}
