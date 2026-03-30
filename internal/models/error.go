package models

type CCError struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	StatusCode   int    `json:"status_code"`
}
