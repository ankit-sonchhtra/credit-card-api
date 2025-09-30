package utils

import (
	"github.com/credit-card-service/internal/models"
	"github.com/credit-card-service/pkg/constants"
)

func NewCCInternalServerError() *models.CCError {
	return &models.CCError{
		ErrorCode:      constants.InternalServerErrCode,
		ErrorMessage:   constants.InternalServerErrMsg,
		AdditionalData: models.AdditionalData{StatusCode: 500},
	}
}

func NewCCBadRequestError(errorMsg string) *models.CCError {
	return &models.CCError{
		ErrorCode:      constants.BadRequestErrCode,
		ErrorMessage:   errorMsg,
		AdditionalData: models.AdditionalData{StatusCode: 400},
	}
}
