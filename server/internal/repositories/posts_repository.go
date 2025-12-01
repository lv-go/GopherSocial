package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
)

type PostsRepository struct {
	gormCRUDRepository[models.Post, uint]
}

func NewPostsRepository() *PostsRepository {
	return &PostsRepository{
		gormCRUDRepository[models.Post, uint]{
			db: gormDB,
		},
	}
}

func (repo *PostsRepository) GetPageByUserID(ctx context.Context, userID string, page PageQuery) (*Page[models.Post], error) {
	return repo.GetPage(ctx, models.Post{UserID: userID}, page)
}
