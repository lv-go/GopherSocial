package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type UsersService struct {
	//usersRepository           repositories.UsersRepository
	//userInvitationsRepository repositories.UserInvitationsRepository
	db *gorm.DB
}

func (r *UsersService) Activate(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. find the user that this token belongs to
		user, err := r.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. update the user
		user.IsActive = true
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		// 3. clean the invitations
		return tx.Delete(models.UserInvitation{
			UserID: user.ID,
		}).Error
	})
}

func (r *UsersService) getUserFromInvitation(ctx context.Context, tx *gorm.DB, token string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &models.User{}
	tx.Find(user, query, hashToken, time.Now())
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
