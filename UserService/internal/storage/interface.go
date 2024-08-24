package storage

import "context"

type Storage interface {
	Get(ctx context.Context, key string) (val any, err error)
	Set(ctx context.Context, key string, value any) (err error)
	Delete(ctx context.Context, key string) (err error)
}
