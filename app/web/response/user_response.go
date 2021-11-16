package response

import (
	"github.com/jinzhu/gorm"
	"github.com/zidni722/login-crud-user/app/dto/response"
	"github.com/zidni722/login-crud-user/app/models"
)

type UserResponse struct {
	Db *gorm.DB
}

func NewUserResponse(db *gorm.DB) UserResponse {
	return UserResponse{Db: db}
}

func (r *UserResponse) New(user models.User) response.User {
	response := response.User{
		ID:   user.ID,
		Username: user.Username,
		Email: user.Email,
		Address: user.Address,
		Password: user.Password,
	}

	return response
}

func (r *UserResponse) Collection(users []models.User) []response.User {
	var responses []response.User

	for _, user := range users {
		responses = append(responses, r.New(user))
	}

	return responses
}
