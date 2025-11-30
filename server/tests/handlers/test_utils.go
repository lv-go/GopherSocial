package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sikozonpc/social/internal/app"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/handlers"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/repositories"
	"go.uber.org/zap"
)

func setupTestApplication(t *testing.T, cfg config.Config) {
	t.Helper()

	app.Logger = zap.NewNop().Sugar()
	// Uncomment to enable logs
	// logger := zap.Must(zap.NewProduction()).Sugar()

	// Rate limiter
	app.RateLimiter = ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	app.Config = cfg

	app.RateLimiterMiddlewares = ratelimiter.NewMiddlewares(cfg.RateLimiter, app.RateLimiter)

	app.UsersRepository = repositories.NewUsersRepository()

	app.AuthMiddlewares = auth.NewMiddlewares(
		app.UsersRepository,
		app.RateLimiter,
		app.AuthClient,
		app.Logger,
		app.Config,
	)
	app.UsersHandlers = handlers.NewUsersHandlers(app.AuthMiddlewares, app.UsersRepository, app.FollowersRepository)
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
