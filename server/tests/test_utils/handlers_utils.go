package test_utils

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sikozonpc/social/internal/app"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/handlers"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/repositories"
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
	user1 := GetUser1()
	authClient := new(mockAuthClient)
	authClient.On("VerifyIDToken", mock.Anything, "invalid_token").
		Return(nil, errors.New("invalid_token"))
	authClient.On("VerifyIDToken", mock.Anything, "valid_token").
		Return(&auth.Token{
			UID: user1.ID,
			Claims: map[string]interface{}{
				"email":          user1.Email,
				"email_verified": user1.IsActive,
			},
		})

	app.AuthClient = authClient
}

func SetupTestApplication(t *testing.T) {
	t.Helper()
	cfg := app.Config

	app.Logger = zap.NewNop().Sugar()
	// Uncomment to enable logs
	// logger := zap.Must(zap.NewProduction()).Sugar()

	// Rate limiter
	app.RateLimiter = ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	setupMockAuthClient()
	app.PostRepository = repositories.NewPostsRepository()
	app.CommentsRepository = repositories.NewCommentsRepository()
	app.RolesRepository = repositories.NewRolesRepository()

	app.RateLimiterMiddlewares = ratelimiter.NewMiddlewares(cfg.RateLimiter, app.RateLimiter)

	app.AuthMiddlewares = auth.NewMiddlewares(
		app.RateLimiter,
		app.AuthClient,
		app.Logger,
		app.Config,
	)

	app.FeedHandlers = handlers.NewFeedHandlers(app.PostRepository)
	app.PostsHandlers = handlers.NewPostsHandlers(app.PostRepository, app.CommentsRepository, app.RolesRepository)
}

func ExecuteRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
