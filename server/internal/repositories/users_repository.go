package repositories

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type UsersRepository struct {
	gormCRUDRepository  CRUDRepository[models.User, uuid.UUID]
	redisCRUDRepository CRUDRepository[models.User, uuid.UUID]
	db                  *gorm.DB
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{
		gormCRUDRepository:  NewGormCRUDRepository[models.User, uuid.UUID](),
		redisCRUDRepository: NewRedisCRUDRepository[models.User, uuid.UUID]("user", time.Minute),
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

func (r *UsersRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := r.redisCRUDRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		slog.Debug("User found in Redis cache", "user", user)
		return user, nil
	}
	user, err = r.gormCRUDRepository.GetByID(ctx, id)
	err = r.redisCRUDRepository.Create(ctx, user)
	if err != nil {
		slog.Error("Error adding user to Redis cache", "error", err)
		return nil, err
	}
	slog.Debug("User found in DB", "user", user)
	return user, nil
}

func (r *UsersRepository) GetOne(ctx context.Context, filter interface{}) (*models.User, error) {
	user, err := r.redisCRUDRepository.GetOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	if user != nil {
		slog.Debug("User found in Redis cache", "user", user)
		return user, nil
	}
	user, err = r.gormCRUDRepository.GetOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = r.redisCRUDRepository.Create(ctx, user)
	if err != nil {
		slog.Error("Error adding user to Redis cache", "error", err)
		return nil, err
	}
	slog.Debug("User found in DB", "user", user)
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

func (r *UsersRepository) UpdateByID(ctx context.Context, id uuid.UUID, entity *models.User) error {
	err := r.gormCRUDRepository.UpdateByID(ctx, id, entity)
	if err != nil {
		return err
	}
	return r.redisCRUDRepository.UpdateByID(ctx, id, entity)
}

func (r *UsersRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
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
		user, err := r.getUserFromInvitation(tx, token)
		if err != nil {
			return err
		}

		// 2. update the user
		if err := tx.Model(user).UpdateColumn("is_active", true).Error; err != nil {
			return err
		}

		// 3. clean the invitations
		return tx.Delete(models.UserInvitation{}, "user_id = ?", user.ID).Error
	})
}

func (r *UsersRepository) getUserFromInvitation(tx *gorm.DB, token string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active, u.role_id
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &models.User{}
	tx.Raw(query, hashToken, time.Now()).Scan(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (r *UsersRepository) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			return err
		}

		db := tx.WithContext(ctx).Create(models.UserInvitation{
			Token:  token,
			UserID: user.ID,
			Expiry: time.Now().Add(invitationExp),
		})
		if db.Error != nil {
			return db.Error
		}

		return nil
	})
}
