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

func TestCommentsRepository(t *testing.T) {
	testUser := test_utils.InitTestUser()
	testPost := test_utils.InitTestPost(testUser)

	testVal := models.Comment{
		Content: "test comment",
		UserID:  testUser.ID,
		PostID:  testPost.ID,
	}

	expected := models.Comment{
		Content: testVal.Content,
		UserID:  testVal.UserID,
		PostID:  testVal.PostID,
	}

	repository := repositories.NewGormCRUDRepository[models.Comment, uint]()

	t.Run("should create a comment", func(t *testing.T) {
		err := repository.Create(t.Context(), &testVal)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(
			expected,
			testVal,
			cmpopts.IgnoreFields(
				models.Comment{},
				"ID", "CreatedAt", "UpdatedAt",
			),
		)
		assert.Empty(t, diff, diff)
		assert.NotEmpty(t, testVal.ID)
		expected.ID = testVal.ID
		expected.CreatedAt = testVal.CreatedAt
		expected.UpdatedAt = testVal.UpdatedAt
	})
	t.Run("should get comment by id", func(t *testing.T) {
		actual, err := repository.GetByID(t.Context(), testVal.ID)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(*actual, testVal, cmpopts.EquateApproxTime(time.Minute))
		assert.Empty(t, diff, diff)
	})
	t.Run("should update comment by id", func(t *testing.T) {
		testVal.Content = "updated test comment"
		expected.Content = testVal.Content
		err := repository.UpdateByID(t.Context(), testVal.ID, &testVal)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		diff := cmp.Diff(expected, testVal,
			cmpopts.IgnoreFields(models.Comment{}, "UpdatedAt"),
		)
		assert.Empty(t, diff, diff)
		assert.Greater(t, testVal.UpdatedAt, expected.UpdatedAt)
	})
	t.Run("should delete comment by id", func(t *testing.T) {
		err := repository.DeleteByID(t.Context(), testVal.ID)
		if err != nil {
			t.Fatalf("should not error, but got %v", err)
		}
		actual, err := repository.GetByID(t.Context(), testVal.ID)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, actual)
	})
}
