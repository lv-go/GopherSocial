package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sikozonpc/social/internal/api"
	"github.com/sikozonpc/social/internal/handlers"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/tests/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestPostsHandlers(t *testing.T) {
	user1 := test_utils.GetUser1()
	test_utils.SetupTestApplication(t)
	mux := api.Mount()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	expected := models.Post{
		Content:   "test post",
		Title:     "test post",
		UserID:    user1.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	t.Run("should return unauthorized error when no token is provided", func(t *testing.T) {
		req := httptest.NewRequest("POST", ts.URL+"/v1/posts", nil)
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 401, rr.Code)
	})
	t.Run("should return bad request error when valid token is provided but invalid payload", func(t *testing.T) {
		req := httptest.NewRequest("POST", ts.URL+"/v1/posts", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 400, rr.Code)
	})
	t.Run("should create a post when valid token is provided", func(t *testing.T) {
		createPostPayload := handlers.CreatePostPayload{
			Title:   expected.Title,
			Content: expected.Content,
		}
		jsonTestValue, _ := json.Marshal(createPostPayload)
		req := httptest.NewRequest("POST", ts.URL+"/v1/posts", bytes.NewBuffer(jsonTestValue))
		req.Header.Set("Authorization", "Bearer valid_token")
		req.Header.Set("Content-Type", "application/json")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 201, rr.Code)

		var actual models.Post
		err := json.NewDecoder(rr.Body).Decode(&actual)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(
			expected,
			actual,
			cmpopts.IgnoreFields(models.Post{}, "ID"),
			cmpopts.EquateApproxTime(time.Minute),
		)
		assert.Empty(t, diff)
		assert.NotEmpty(t, actual.ID)
		expected.ID = actual.ID
	})
	t.Run("should get post by id when valid token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("%s/v1/posts/%d", ts.URL, expected.ID), nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 200, rr.Code)

		var actual models.Post
		err := json.NewDecoder(rr.Body).Decode(&actual)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}

		expected.Comments = []models.Comment{}

		diff := cmp.Diff(
			expected,
			actual,
			cmpopts.EquateApproxTime(time.Minute),
		)
		assert.Empty(t, diff)
	})
	t.Run("should update post by id when valid token is provided", func(t *testing.T) {
		expected.Title = "updated test post"
		expected.Content = "updated test post"

		updatePostPayload := handlers.UpdatePostPayload{
			Title:   &expected.Title,
			Content: &expected.Content,
		}
		jsonTestValue, _ := json.Marshal(updatePostPayload)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("%s/v1/posts/%d", ts.URL, expected.ID),
			bytes.NewBuffer(jsonTestValue))

		req.Header.Set("Authorization", "Bearer valid_token")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 200, rr.Code)

		var actual models.Post
		err := json.NewDecoder(rr.Body).Decode(&actual)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}

		expected.Comments = nil

		diff := cmp.Diff(
			expected,
			actual,
			cmpopts.EquateApproxTime(time.Minute),
		)
		assert.Empty(t, diff)
	})
	t.Run("should delete post by id when valid token is provided", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("%s/v1/posts/%d", ts.URL, expected.ID), nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, 204, rr.Code)
	})
}
