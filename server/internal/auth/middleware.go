package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/utils"
	"go.uber.org/zap"
)

type Middlewares struct {
	rateLimiter ratelimiter.RateLimiter
	logger      *zap.SugaredLogger
	config      config.Config
	authClient  Client
}

func NewMiddlewares(
	rateLimiter ratelimiter.RateLimiter,
	authClient Client,
	logger *zap.SugaredLogger,
	config config.Config,
) Middlewares {
	return Middlewares{
		rateLimiter: rateLimiter,
		authClient:  authClient,
		logger:      logger,
		config:      config,
	}
}

func (m *Middlewares) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		idToken := parts[1]
		jwtToken, err := m.authClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			utils.UnauthorizedErrorResponse(w, r, err)
			return
		}

		// The UID from the token is not a standard UUID. We generate a deterministic UUIDv5 from it.
		userId := uuid.NewSHA1(uuid.NameSpaceURL, []byte(jwtToken.UID))
		r = SetUserInContext(r, &models.User{
			ID:       userId,
			Email:    jwtToken.Claims["email"].(string),
			IsActive: jwtToken.Claims["email_verified"].(bool),
		})
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				utils.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				utils.UnauthorizedBasicErrorResponse(w, r, err)
				return
			}

			// check the credentials
			username := m.config.Auth.Basic.User
			pass := m.config.Auth.Basic.Pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				utils.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middlewares) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.RateLimiter.Enabled {
			if allow, retryAfter := m.rateLimiter.Allow(r.RemoteAddr); !allow {
				utils.RateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
