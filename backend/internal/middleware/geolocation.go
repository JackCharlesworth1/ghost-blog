package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
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
	log.Printf("GetCountryFromIP - Input IP: %s", ipAddress)

	// Handle localhost and private IPs
	if ipAddress == "127.0.0.1" || ipAddress == "::1" || ipAddress == "" {
		log.Printf("GetCountryFromIP - Detected as localhost/empty")
		return "Local"
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		log.Printf("GetCountryFromIP - Failed to parse IP")
		return "Unknown"
	}

	if ip.IsPrivate() || ip.IsLoopback() {
		log.Printf("GetCountryFromIP - IP is private or loopback (IsPrivate: %v, IsLoopback: %v)", ip.IsPrivate(), ip.IsLoopback())
		return "Local"
	}

	log.Printf("GetCountryFromIP - IP is public, querying ip-api.com")

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
	// Format: "client-ip, proxy1-ip, proxy2-ip" - first IP is the original client
	forwarded := r.Header.Get("X-Forwarded-For")
	log.Printf("GetClientIP - X-Forwarded-For: %s", forwarded)
	if forwarded != "" {
		// Split by comma and get the first IP (original client)
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	log.Printf("GetClientIP - X-Real-IP: %s", realIP)
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	log.Printf("GetClientIP - RemoteAddr: %s", r.RemoteAddr)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
