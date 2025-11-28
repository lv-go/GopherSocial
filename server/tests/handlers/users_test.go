package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sikozonpc/social/internal/api"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/tests/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	cfg := config.Config{
		RedisCfg: repositories.RedisConfig{
			Enabled: true,
		},
		GormDBConfig: repositories.GormDBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "admin",
			Password: "S3cureP@ssw0rd!",
			Dbname:   "gopher_social",
		},
		Auth: config.AuthConfig{
			Token: config.TokenConfig{
				Secret: "auth-secret",
				Exp:    72 * time.Hour,
				Iss:    "gophersocial",
			},
		},
	}
	test_utils.InitTestUser()

	setupTestApplication(t, cfg)
	mux := api.Mount()

	var testToken string

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should login the user", func(t *testing.T) {
		jsonBody := `{"email": "testuser@email.com", "password": "testPassword"}`
		req, err := http.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(jsonBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		var loginResponse auth.TokenResponse
		err = json.Unmarshal(rr.Body.Bytes(), &loginResponse)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, loginResponse.Token)
		testToken = loginResponse.Token
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})

	//t.Run("should hit the cache first and if not exists it sets the User on the cache", func(t *testing.T) {
	//	mockCacheStore := app.CacheStorage.Users.(*cache.MockUserStore)
	//
	//	mockCacheStore.On("Get", int64(42)).Return(nil, nil)
	//	mockCacheStore.On("Get", int64(1)).Return(nil, nil)
	//	mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)
	//
	//	req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	req.Header.Set("Authorization", "Bearer "+testToken)
	//
	//	rr := executeRequest(req, mux)
	//
	//	checkResponseCode(t, http.StatusOK, rr.Code)
	//
	//	mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
	//
	//	mockCacheStore.Calls = nil // Reset mock expectations
	//})
	//
	//t.Run("should NOT hit the cache if it is not enabled", func(t *testing.T) {
	//	req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	req.Header.Set("Authorization", "Bearer "+testToken)
	//
	//	rr := executeRequest(req, mux)
	//
	//	checkResponseCode(t, http.StatusOK, rr.Code)
	//})
}
