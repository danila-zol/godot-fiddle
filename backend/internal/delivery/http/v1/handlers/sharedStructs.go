package handlers

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ObjectUploader interface {
	CheckFileSize(size int64, userTier string) error

	ObjectTooLargeErr() error
	ObjectNotFoundErr() error
}
