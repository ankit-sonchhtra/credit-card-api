package models

import (
	"errors"
	"fmt"

	"github.com/credit-card-api/pkg/constants"
	"github.com/go-playground/validator/v10"
)

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required,max=12,numeric" example:"0987654321"` // only allow numeric value
}

type CreateAccountResponse struct {
	AccountId      int64  `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"0987654321"`
}

type GetAccountResponse struct {
	AccountId      int64  `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"0987654321"`
}

func (request CreateAccountRequest) Validate() error {
	err := validator.New().Struct(&request)
	return translateError(err)
}

func translateError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fe := ve[0]
		field := fe.Field()
		tag := fe.Tag()
		param := fe.Param()
		switch tag {
		case constants.RequiredTag:
			return errors.New(fmt.Sprintf("The '%s' field is mandatory.", field))
		case constants.MaxTag:
			return errors.New(fmt.Sprintf("The '%s' field cannot exceed %s digits.", field, param))
		case constants.NumericTag:
			return errors.New(fmt.Sprintf("The '%s' field will only accept numeric value.", field))
		case constants.GTTag:
			return errors.New(fmt.Sprintf("The '%s' field value must be greater than %s.", field, param))
		}
	}
	return err
}
