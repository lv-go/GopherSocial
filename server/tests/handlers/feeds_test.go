package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/sikozonpc/social/internal/api"
	"github.com/sikozonpc/social/internal/app"
	"github.com/sikozonpc/social/tests/test_utils"
)

func TestFeedsHandlers(t *testing.T) {
	user1 := test_utils.GetUser1()
	test_utils.InitTestPost(user1)
	app.Setup(t.Context(), "../../")
	setupMockAuthClient()
	mux := api.Mount()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	t.Run("should return unauthorized error when no token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", ts.URL+"/v1/user/feed", nil)
		rr := executeRequest(req, mux)
		checkResponseCode(t, 401, rr.Code)
	})
	t.Run("should return unauthorized error when invalid token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", ts.URL+"/v1/user/feed", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		rr := executeRequest(req, mux)
		checkResponseCode(t, 401, rr.Code)
	})
	t.Run("should return ok when valid token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", ts.URL+"/v1/user/feed", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		rr := executeRequest(req, mux)
		checkResponseCode(t, 200, rr.Code)
	})
}
