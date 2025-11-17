package repositories

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type UsersRepository struct {
	gormCRUDRepository  CRUDRepository[models.User, uint]
	redisCRUDRepository CRUDRepository[models.User, uint]
	db                  *gorm.DB
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{
		gormCRUDRepository:  NewGormCRUDRepository[models.User, uint](),
		redisCRUDRepository: NewRedisCRUDRepository[models.User, uint]("user", time.Minute),
		db:                  gormDB,
	}
}

func (r *UsersRepository) Create(ctx context.Context, user *models.User) error {
	err := r.gormCRUDRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return r.redisCRUDRepository.Create(ctx, user)
}

func (r *UsersRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := r.redisCRUDRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return r.gormCRUDRepository.GetByID(ctx, id)
	}
	return user, nil
}

func (r *UsersRepository) GetOne(ctx context.Context, filter interface{}) (*models.User, error) {
	user, err := r.redisCRUDRepository.GetOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return r.gormCRUDRepository.GetOne(ctx, filter)
	}
	return user, nil
}

func (r *UsersRepository) GetAll(ctx context.Context, filter interface{}) ([]models.User, error) {
	return r.gormCRUDRepository.GetAll(ctx, filter)
}

func (r *UsersRepository) GetPage(
	ctx context.Context,
	filter map[string]interface{},
	page PageQuery,
) (*Page[models.User], error) {
	return r.gormCRUDRepository.GetPage(ctx, filter, page)
}

func (r *UsersRepository) UpdateByID(ctx context.Context, id uint, entity *models.User) error {
	err := r.gormCRUDRepository.UpdateByID(ctx, id, entity)
	if err != nil {
		return err
	}
	return r.redisCRUDRepository.UpdateByID(ctx, id, entity)
}

func (r *UsersRepository) DeleteByID(ctx context.Context, id uint) error {
	err := r.gormCRUDRepository.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return r.redisCRUDRepository.DeleteByID(ctx, id)
}

func (r *UsersRepository) DeleteOne(ctx context.Context, filter interface{}) error {
	err := r.gormCRUDRepository.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return r.redisCRUDRepository.DeleteOne(ctx, filter)
}

func (r *UsersRepository) GetOneByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.GetOne(ctx, models.User{
		Username: username,
	})
}

func (r *UsersRepository) GetOneByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.GetOne(ctx, models.User{
		Email: email,
	})
}

func (r *UsersRepository) Activate(ctx context.Context, token string) error {
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

func (r *UsersRepository) getUserFromInvitation(ctx context.Context, tx *gorm.DB, token string) (*models.User, error) {
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
