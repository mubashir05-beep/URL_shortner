package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// IPInfo struct to hold API response
type IPInfo struct {
	Country string `json:"country"`
}

// GetCountryFromIP fetches country from IP address
func GetCountryFromIP(ip string) string {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	var data IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "Unknown"
	}

	return data.Country
}
