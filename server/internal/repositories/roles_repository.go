package repositories

import (
	"context"

	"github.com/sikozonpc/social/internal/models"
	"gorm.io/gorm"
)

type RolesRepository struct {
	gormCRUDRepository[models.Role, uint]
	db *gorm.DB
}

func NewRolesRepository(db *gorm.DB) *RolesRepository {
	return &RolesRepository{
		db: db,
	}
}

func (r *RolesRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	return r.gormCRUDRepository.GetOne(ctx, models.Role{
		Name: name,
	})
}
