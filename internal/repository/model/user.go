package model

type UserDocument struct {
	UserId       string `bson:"user_id"`
	Name         string `bson:"name"`
	Email        string `bson:"email"`
	MobileNumber string `bson:"mobile_number"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}
