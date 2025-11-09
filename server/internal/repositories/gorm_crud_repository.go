package repositories

import (
	"context"
	"log/slog"
	"reflect"

	"gorm.io/gorm"
)

type gormCRUDRepository[T interface{}, ID any] struct {
	DB gorm.DB
}

func NewGormCRUDRepository[T interface{}, ID any]() CRUDRepository[T, ID] {
	return &gormCRUDRepository[T, ID]{DB: *DB}
}

func (g gormCRUDRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	result := g.DB.WithContext(ctx).Create(entity)
	return result.Error
}

func (g gormCRUDRepository[T, ID]) GetByID(ctx context.Context, id ID) (*T, error) {
	var entity T
	result := g.DB.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (g gormCRUDRepository[T, ID]) GetOne(ctx context.Context, filter interface{}) (*T, error) {
	var entity T
	result := g.DB.WithContext(ctx).Where(filter).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (g gormCRUDRepository[T, ID]) GetAll(ctx context.Context, filter map[string]interface{}) ([]T, error) {
	entities := make([]T, 0)
	result := g.DB.WithContext(ctx).Where(filter).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

func (g gormCRUDRepository[T, ID]) GetPage(ctx context.Context, filter map[string]interface{}, page int, pageSize int) (*Page[T], error) {
	entities := make([]T, 0)
	offset := (page - 1) * pageSize
	result := g.DB.WithContext(ctx).Where(filter).Offset(offset).Limit(pageSize).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	g.DB.WithContext(ctx).Model(new(T)).Where(filter).Count(&total)
	return &Page[T]{
		Items:    entities,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (g gormCRUDRepository[T, ID]) UpdateByID(ctx context.Context, id ID, entity *T) error {
	setId(entity, id)
	result := g.DB.WithContext(ctx).Save(entity)
	return result.Error
}

func (g gormCRUDRepository[T, ID]) DeleteByID(ctx context.Context, id ID) error {
	result := g.DB.WithContext(ctx).Delete(new(T), id)
	return result.Error
}

func (g gormCRUDRepository[T, ID]) DeleteOne(ctx context.Context, filter interface{}) error {
	result := g.DB.WithContext(ctx).Where(filter).Delete(new(T))
	return result.Error
}

func setId[T interface{}, ID any](entity *T, id ID) {
	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		slog.Error("entity must be a pointer")
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		slog.Error("entity must be a struct")
		return
	}
	field := val.FieldByName("ID")
	if !field.IsValid() {
		slog.Error("entity must have an ID field")
		return
	}
	field.Set(reflect.ValueOf(id))
}
