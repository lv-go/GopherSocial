package ratelimiter

import (
	"net/http"

	"github.com/sikozonpc/social/internal/utils"
)

type Middlewares struct {
	config      Config
	rateLimiter RateLimiter
}

func NewMiddlewares(config Config, rateLimiter RateLimiter) Middlewares {
	return Middlewares{
		config:      config,
		rateLimiter: rateLimiter,
	}
}

func (m *Middlewares) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.Enabled {
			if allow, retryAfter := m.rateLimiter.Allow(r.RemoteAddr); !allow {
				utils.RateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
