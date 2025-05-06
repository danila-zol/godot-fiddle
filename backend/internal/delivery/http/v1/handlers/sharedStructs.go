package handlers

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
