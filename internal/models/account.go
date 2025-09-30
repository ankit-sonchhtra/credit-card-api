package models

type CreateAccountRequest struct {
	UserId         string `json:"userId" binding:"required" example:"92d68c0e-dafe-406a-a0f2-8faae2020947"`
	DocumentNumber string `json:"documentNumber" binding:"required" example:"0987654321"`
}

type CreateAccountResponse struct {
	AccountId      string `json:"accountId" example:"92d68c0e-dafe-406a-a0f2-8faae2020947"`
	DocumentNumber string `json:"documentNumber" example:"0987654321"`
}

type GetAccountResponse struct {
	AccountId      string `json:"accountId" example:"92d68c0e-dafe-406a-a0f2-8faae2020947"`
	UserId         string `json:"userId" example:"92d68c0e-dafe-406a-a0f2-8faae2020947"`
	DocumentNumber string `json:"documentNumber" example:"0987654321"`
}
