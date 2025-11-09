package test_utils

import (
	"context"
	"errors"

	"github.com/joho/godotenv"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/repositories"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load("../../.env.local")
	if err != nil {
		panic(err)
	}

	repositories.SetupDB()

}

func InitTestUser() *models.User {
	userRepository := repositories.NewGormCRUDRepository[models.User, uint]()
	testUser, err := userRepository.GetOne(context.Background(), models.User{Username: "testUser"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			testUser = &models.User{
				Username: "testUser",
				Email:    "testuser@email.com",
				Password: "testPassword",
				IsActive: true,
				RoleID:   1,
			}
			err = userRepository.Create(context.Background(), testUser)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return testUser
}

func InitTestPost(testUser *models.User) *models.Post {
	postRepository := repositories.NewGormCRUDRepository[models.Post, uint]()
	testPost, err := postRepository.GetOne(context.Background(), models.Post{Title: "Test Post"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			testPost = &models.Post{
				Title:   "Test Post",
				Content: "Test Post",
				UserID:  testUser.ID,
			}
			err = postRepository.Create(context.Background(), testPost)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return testPost
}
