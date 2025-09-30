package models

type CCError struct {
	ErrorCode      string      `json:"errorCode"`
	ErrorMessage   string      `json:"errorMessage"`
	AdditionalData interface{} `json:"additionalData,omitempty"`
}

type AdditionalData struct {
	StatusCode int `json:"statusCode"`
}
