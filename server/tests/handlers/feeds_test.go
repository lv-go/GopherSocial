package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/sikozonpc/social/internal/api"
	"github.com/sikozonpc/social/tests/test_utils"
)

func TestFeedsHandlers(t *testing.T) {
	user1 := test_utils.GetUser1()
	test_utils.InitTestPost(user1)
	test_utils.SetupTestApplication(t)
	mux := api.Mount()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	t.Run("should return unauthorized error when no token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", ts.URL+"/v1/user/feed", nil)
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 401, rr.Code)
	})
	t.Run("should return unauthorized error when invalid token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", ts.URL+"/v1/user/feed", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 401, rr.Code)
	})
	t.Run("should return one post when valid token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", ts.URL+"/v1/user/feed", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 200, rr.Code)
	})
}
