package response

import (
	"github.com/zidni722/login-crud-user/app/dto/response"
	"github.com/jinzhu/gorm"
)

type LoginResponse struct {
	Db *gorm.DB
}

func NewLoginResponse(db *gorm.DB) LoginResponse {
	return LoginResponse{Db: db}
}

func (r *LoginResponse) New(token string) response.Login {
	response := response.Login{Token: token}
	return response
}
