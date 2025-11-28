package repositories

import "context"

type CRUDRepository[T interface{}, ID any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id ID) (*T, error)
	GetOne(ctx context.Context, filter interface{}) (*T, error)
	GetAll(ctx context.Context, filter interface{}) ([]T, error)
	GetPage(ctx context.Context, filter interface{}, pageQuery PageQuery) (*Page[T], error)
	UpdateByID(ctx context.Context, id ID, entity *T) error
	DeleteByID(ctx context.Context, id ID) error
	DeleteOne(ctx context.Context, filter interface{}) error
}

type Page[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"`
	Number int   `json:"number"`
	Size   int   `json:"size"`
}

type PageQuery struct {
	Size   int    `json:"size" validate:"gte=1,lte=1000"`
	Number int    `json:"number" validate:"gte=0"`
	Sort   []Sort `json:"sort"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order" validate:"oneof=asc desc"`
}
