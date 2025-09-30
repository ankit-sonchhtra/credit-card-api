package models

type CreateUserRequest struct {
	Name         string `json:"name" example:"John Deo"`
	Email        string `json:"email" binding:"omitempty,email" example:"john.deo@xyz.com"`
	MobileNumber string `json:"mobileNumber" binding:"required" example:"+919825212345"`
}

type CreateUserResponse struct {
	UserId string `json:"userId" example:"92d68c0e-dafe-406a-a0f2-8faae2020947"`
}
