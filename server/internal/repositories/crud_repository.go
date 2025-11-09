package repositories

import "context"

type CRUDRepository[T interface{}, ID any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id ID) (*T, error)
	GetOne(ctx context.Context, filter interface{}) (*T, error)
	GetAll(ctx context.Context, filter map[string]interface{}) ([]T, error)
	GetPage(ctx context.Context, filter map[string]interface{}, page int, pageSize int) (*Page[T], error)
	UpdateByID(ctx context.Context, id ID, entity *T) error
	DeleteByID(ctx context.Context, id ID) error
	DeleteOne(ctx context.Context, filter interface{}) error
}

type Page[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}
