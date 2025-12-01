package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type UserInvitationsRepository struct {
	gormCRUDRepository[models.UserInvitation, uint]
	db *gorm.DB
}

func NewUserInvitationsRepository(
	db *gorm.DB,
) *UserInvitationsRepository {
	return &UserInvitationsRepository{
		db: db,
	}
}

func (r *UserInvitationsRepository) DeleteByUserID(ctx context.Context, userId string) error {
	return r.DeleteOne(ctx, models.UserInvitation{UserID: userId})
}
