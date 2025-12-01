package test_utils

import (
	"context"
	"errors"

	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/repositories"
	"gorm.io/gorm"
)

func init() {
	cfg := config.Setup("../../")
	repositories.SetupGormDB(cfg.GormDBConfig)
	repositories.SetupRedisClient(cfg.RedisCfg)
}

func GetUser1() *models.User {
	return &models.User{
		ID:       "GhTjK7pQ2xR5wVfB9zLdYc4m",
		Email:    "user1@email.com",
		IsActive: true,
	}
}

func InitTestPost(testUser *models.User) *models.Post {
	postRepository := repositories.NewPostsRepository()
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
