package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type GeoResponse struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Status      string `json:"status"`
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func GetCountryFromIP(ipAddress string) string {
	// Handle localhost and private IPs
	if ipAddress == "127.0.0.1" || ipAddress == "::1" || ipAddress == "" {
		return "Local"
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "Unknown"
	}

	if ip.IsPrivate() || ip.IsLoopback() {
		return "Local"
	}

	// Query ip-api.com for geolocation
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode", ipAddress)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Unknown"
	}

	var geoResp GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return "Unknown"
	}

	if geoResp.Status != "success" {
		return "Unknown"
	}

	return geoResp.Country
}

func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
