package utils

import "github.com/mssola/useragent"

// ParseUserAgent extracts browser and device from User-Agent string
func ParseUserAgent(uaString string) (browser, device string) {
	ua := useragent.New(uaString)
	browser, _ = ua.Browser()
	device = "Desktop"
	if ua.Mobile() {
		device = "Mobile"
	} else if ua.Bot() {
		device = "Bot"
	}
	return browser, device
}
