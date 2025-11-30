package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sikozonpc/social/internal/app"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockAuthClient struct {
	mock.Mock
}

func (m *mockAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Token), nil
}

func setupMockAuthClient() {
	authClient := new(mockAuthClient)
	authClient.On("VerifyIDToken", mock.Anything, "invalid_token").
		Return(nil, errors.New("invalid_token"))
	authClient.On("VerifyIDToken", mock.Anything, "valid_token").
		Return(&auth.Token{
			UID: "00000001-0000-0000-0000-000000000001",
			Claims: map[string]interface{}{
				"email":          "user1@email.com",
				"email_verified": true,
			},
		})

	app.AuthClient = authClient
}

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

	app.AuthMiddlewares = auth.NewMiddlewares(
		app.RateLimiter,
		app.AuthClient,
		app.Logger,
		app.Config,
	)
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
