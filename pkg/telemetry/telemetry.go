package telemetry

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var telemetryHTTPClient = &http.Client{Timeout: 5 * time.Second}

func telemetryURL() string {
	return strings.TrimSpace(os.Getenv("TELEMETRY_URL"))
}

// isTelemetryEnabled is true only when TELEMETRY_ENABLED=true and TELEMETRY_URL is non-empty.
func isTelemetryEnabled() bool {
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("TELEMETRY_ENABLED")), "true") {
		return false
	}
	return telemetryURL() != ""
}

type TelemetryData struct {
	Route      string    `json:"route"`
	APIVersion string    `json:"apiVersion"`
	Timestamp  time.Time `json:"timestamp"`
}

type telemetryService struct{}

func (t *telemetryService) TelemetryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isTelemetryEnabled() {
			c.Next()
			return
		}
		route := c.FullPath()
		go SendTelemetry(route)
		c.Next()
	}
}

type TelemetryService interface {
	TelemetryMiddleware() gin.HandlerFunc
}

func SendTelemetry(route string) {
	if !isTelemetryEnabled() || route == "/" {
		return
	}

	telemetry := TelemetryData{
		Route:      route,
		APIVersion: "evo-go",
		Timestamp:  time.Now(),
	}

	url := telemetryURL()
	data, err := json.Marshal(telemetry)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := telemetryHTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func NewTelemetryService() TelemetryService {
	return &telemetryService{}
}
