package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
)

type CommentsRepository struct {
	gormCRUDRepository[models.Comment, uint]
}

func NewCommentsRepository() *CommentsRepository {
	return &CommentsRepository{
		gormCRUDRepository: gormCRUDRepository[models.Comment, uint]{
			db: gormDB,
		},
	}
}

func (r *CommentsRepository) GetByPostID(ctx context.Context, postID uint) ([]models.Comment, error) {
	return r.GetAll(ctx, models.Comment{
		PostID: postID,
	})
}
