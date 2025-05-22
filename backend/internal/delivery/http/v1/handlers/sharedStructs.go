package handlers

import (
	// "context"
	"io"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ObjectUploader interface {
	CheckFileSize(size int64, userTier string) error

	PutObject(objectKey string, file io.Reader) error
	// GetObject(ctx context.Context, objectKey string, fileName string) (*[]byte, error)
	GetObjectLink(objectKey string) (string, error)
	DeleteObject(objectKey string) error

	ObjectTooLargeErr() error
	ObjectNotFoundErr() error
}
