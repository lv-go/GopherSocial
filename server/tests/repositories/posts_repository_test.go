package repositories

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/tests/test_utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostsRepository(t *testing.T) {
	user1 := test_utils.GetUser1()
	testVal := models.Post{
		Title:   "test post",
		Content: "test post",
		UserID:  user1.ID,
	}

	expected := models.Post{
		Title:   testVal.Title,
		Content: testVal.Content,
		UserID:  testVal.UserID,
	}

	repository := repositories.NewPostsRepository()

	t.Run("should create a post", func(t *testing.T) {
		err := repository.Create(t.Context(), &testVal)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(
			expected,
			testVal,
			cmpopts.IgnoreFields(
				models.Post{},
				"ID", "CreatedAt", "UpdatedAt",
			),
		)
		assert.Empty(t, diff, diff)
		assert.NotEmpty(t, testVal.ID)
		expected.ID = testVal.ID
		expected.CreatedAt = testVal.CreatedAt
		expected.UpdatedAt = testVal.UpdatedAt
	})
	t.Run("should get post by id", func(t *testing.T) {
		actual, err := repository.GetByID(t.Context(), testVal.ID)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(*actual, testVal, cmpopts.EquateApproxTime(time.Minute))
		assert.Empty(t, diff, diff)
	})
	t.Run("should update post by id", func(t *testing.T) {
		testVal.Content = "updated test post"
		expected.Content = testVal.Content
		err := repository.UpdateByID(t.Context(), testVal.ID, &testVal)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(expected, testVal,
			cmpopts.IgnoreFields(models.Post{}, "UpdatedAt"),
		)
		assert.Empty(t, diff, diff)
		assert.Greater(t, testVal.UpdatedAt, expected.UpdatedAt)
	})
	t.Run("should delete post by id", func(t *testing.T) {
		err := repository.DeleteByID(t.Context(), testVal.ID)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		actual, err := repository.GetByID(t.Context(), testVal.ID)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, actual)
	})
}
