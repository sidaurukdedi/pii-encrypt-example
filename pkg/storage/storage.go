package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	PutObject(ctx context.Context, bucketName string, filepath string, file io.Reader, contentType string, meta map[string]string) (err error)
	SignURL(ctx context.Context, method, bucketName, filepath string, expiresIn time.Duration) (url string, err error)
}
