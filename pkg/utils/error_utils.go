package utils

import (
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/pkg/constants"
)

func NewCCBadRequestError(errorMsg string) *models.CCError {
	return &models.CCError{
		ErrorCode:    constants.BadRequestErrCode,
		ErrorMessage: errorMsg,
		StatusCode:   400,
	}
}
