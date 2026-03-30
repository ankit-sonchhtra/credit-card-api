package models

type BadRequestError struct {
	ErrorCode    string `json:"error_code" example:"ERR_CC_BAD_REQUEST"`
	ErrorMessage string `json:"error_message" example:"invalid request body"`
	StatusCode   int64  `json:"status_code" example:"400"`
}

type ConflictError struct {
	ErrorCode    string `json:"error_code" example:"ERR_CC_CONFLICT"`
	ErrorMessage string `json:"error_message" example:"resource already exist"`
	StatusCode   int64  `json:"status_code" example:"409"`
}

type InternalServerError struct {
	ErrorCode    string `json:"error_code" example:"ERR_CC_INTERNAL_SERVER_ERROR"`
	ErrorMessage string `json:"error_message" example:"internal server error"`
	StatusCode   int64  `json:"status_code" example:"500"`
}

type NotFoundError struct {
	ErrorCode    string `json:"error_code" example:"ERR_CC_ACCOUNT_NOT_PRESENT"`
	ErrorMessage string `json:"error_message" example:"account not present"`
	StatusCode   int    `json:"status_code" example:"404"`
}

type UnprocessableEntityError struct {
	ErrorCode    string `json:"error_code" example:"ERR_CC_TRANSACTION_ACCOUNT_NOT_FOUND"`
	ErrorMessage string `json:"error_message" example:"account does not exist with provided id."`
	StatusCode   int64  `json:"status_code" example:"422"`
}
