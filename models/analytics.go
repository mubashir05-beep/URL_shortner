package models

import "time"

type Analytics struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URLID     uint      `gorm:"index" json:"url_id"` // Foreign Key for URL
	IPAddress string    `json:"ip_address"`
	Country   string    `json:"country"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	ClickedAt time.Time `json:"clicked_at"`
}
