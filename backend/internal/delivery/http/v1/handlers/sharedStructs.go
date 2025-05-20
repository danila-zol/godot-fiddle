package handlers

import "io"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ObjectUploader interface {
	CheckFileSize(size int64, userTier string) error
	PutObject(objectKey string, file io.Reader) (string, error)
	DeleteObject(objectKey string) error
	ObjectTooLargeErr() error
	ObjectNotFoundErr() error
}
