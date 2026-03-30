package constants

var (
	AccountIdPathParam = "accountId"

	BadRequestErrCode                 = "ERR_CC_BAD_REQUEST"
	InternalServerErrCode             = "ERR_CC_INTERNAL_SERVER_ERROR"
	AccountAlreadyExistErrCode        = "ERR_CC_ACCOUNT_ALREADY_EXIST"
	AccountNotFoundErrCode            = "ERR_CC_ACCOUNT_NOT_FOUND"
	InvalidOperationTypeErrCode       = "ERR_CC_INVALID_OPERATION_TYPE"
	TransactionAccountNotFoundErrCode = "ERR_CC_TRANSACTION_ACCOUNT_NOT_FOUND"

	InvalidRequestBodyErrMsg = "invalid request body"
	AccountIdMissingErrMsg   = "accountId is missing in path params"

	DBUrl = "DB_URL"

	RequiredTag = "required"
	MaxTag      = "max"
	NumericTag  = "numeric"
	GTTag       = "gt"

	EmptyString = ""
)
