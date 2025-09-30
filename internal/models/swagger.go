package models

type BadRequestError struct {
	ErrorCode      string `json:"errorCode" example:"ERR_CC_BAD_REQUEST"`
	ErrorMessage   string `json:"errorMessage" example:"invalid request body"`
	AdditionalData struct {
		StatusCode int64 `json:"statusCode" example:"400"`
	} `json:"additionalData"`
}

type InternalServerError struct {
	ErrorCode      string `json:"errorCode" example:"ERR_CC_INTERNAL_SERVER_ERROR"`
	ErrorMessage   string `json:"errorMessage" example:"internal server error"`
	AdditionalData struct {
		StatusCode int64 `json:"statusCode" example:"500"`
	} `json:"additionalData"`
}

type NotFoundError struct {
	ErrorCode      string `json:"errorCode" example:"ERR_CC_ACCOUNT_NOT_PRESENT"`
	ErrorMessage   string `json:"errorMessage" example:"account not present"`
	AdditionalData struct {
		StatusCode int `json:"statusCode" example:"404"`
	} `json:"additionalData"`
}
