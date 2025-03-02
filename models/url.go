package models

import "time"

type URL struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"` // Foreign Key for User
	ShortCode   string    `gorm:"uniqueIndex" json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int       `gorm:"default:0" json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`

	// Relation (One URL -> Many Analytics)
	Analytics []Analytics `gorm:"foreignKey:URLID" json:"analytics"`
}
