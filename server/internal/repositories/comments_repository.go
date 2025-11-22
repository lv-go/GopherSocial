package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type CommentsRepository struct {
	gormCRUDRepository[models.Comment, uint]
	db *gorm.DB
}

func NewCommentsRepository() *CommentsRepository {
	return &CommentsRepository{
		db: gormDB,
	}
}

func (r *CommentsRepository) GetByPostID(ctx context.Context, postID uint) ([]models.Comment, error) {
	return r.GetAll(ctx, models.Comment{
		PostID: postID,
	})
}
