package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sikozonpc/social/internal/api"
	"github.com/sikozonpc/social/internal/app"
)

func TestRateLimiterMiddleware(t *testing.T) {
	app.Setup(t.Context())
	cfg := app.Config
	ts := httptest.NewServer(api.Mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < cfg.RateLimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		defer resp.Body.Close()

		if i < cfg.RateLimiter.RequestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status Too Many Requests; got %v", resp.Status)
			}
		}
	}
}
