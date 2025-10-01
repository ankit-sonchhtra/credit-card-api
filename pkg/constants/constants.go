package constants

var (
	MongoDatabaseName     = "credit-card-api"
	AccountCollection     = "accounts"
	UserCollection        = "users"
	TransactionCollection = "transactions"

	EmptyString = ""

	AccountIdPathParam = "accountId"
	AccountIdFilter    = "account_id"
	UserIdFilter       = "user_id"
	MobileNumberFilter = "mobile_number"

	BadRequestErrCode        = "ERR_CC_BAD_REQUEST"
	InternalServerErrCode    = "ERR_CC_INTERNAL_SERVER_ERROR"
	AccountNotPresentErrCode = "ERR_CC_ACCOUNT_NOT_PRESENT"
	UserAlreadyExistErrCode  = "ERR_CC_USER_ALREADY_EXIST"

	InvalidRequestBodyErrMsg   = "invalid request body"
	InvalidMobileNumberErrMsg  = "invalid mobile number"
	InternalServerErrMsg       = "internal server error"
	AccountNotPresentErrMsg    = "account not present"
	AmountMustBeNegativeErrMsg = "amount must be negative for purchases and withdrawals"
	AmountMustBePositiveErrMsg = "amount must be positive for payments"
	InvalidOperationTypeErrMsg = "invalid operation type"
	InvalidAccountIdErrMsg     = "account does not exists with requested accountId"
	InvalidUserIdErrMsg        = "user does not exists with requested userId"
	AccountIdMissingErrMsg     = "accountId is missing in path params"
	UserAlreadyExistErrMsg     = "user already exist with requested mobile number"

	OpTypeCashPurchase        = "cash purchase"
	OpTypeInstallmentPurchase = "installment purchase"
	OpTypeWithdrawal          = "withdrawal"
	OpTypePayment             = "payment"

	AsiaKolkataTimeZone = "Asia/Kolkata"

	MongoUri = "MONGO_URI"
)
