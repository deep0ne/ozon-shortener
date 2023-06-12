package storage

import "context"

type DB interface {
	AddLink(ctx context.Context, id int64, url string) (string, error)
	GetLink(ctx context.Context, id int64) (string, error)
}
