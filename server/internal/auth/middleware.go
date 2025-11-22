package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/utils"
	"go.uber.org/zap"
)

type Middlewares struct {
	usersRepository *repositories.UsersRepository
	rateLimiter     ratelimiter.RateLimiter
	authenticator   Authenticator
	logger          *zap.SugaredLogger
	config          config.Config
}

func NewMiddlewares(
	usersRepository *repositories.UsersRepository,
	rateLimiter ratelimiter.RateLimiter,
	authenticator Authenticator,
	logger *zap.SugaredLogger,
	config config.Config,
) Middlewares {
	return Middlewares{
		usersRepository: usersRepository,
		rateLimiter:     rateLimiter,
		authenticator:   authenticator,
		logger:          logger,
		config:          config,
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

		token := parts[1]
		jwtToken, err := m.authenticator.ValidateToken(token)
		if err != nil {
			utils.UnauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			utils.UnauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := m.usersRepository.GetByID(ctx, uint(userID))
		if err != nil {
			utils.UnauthorizedErrorResponse(w, r, err)
			return
		}

		r = SetUserInContext(r, user)
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
