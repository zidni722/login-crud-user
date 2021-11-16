package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/zidni722/login-crud-user/app/dto/request"
	_interface "github.com/zidni722/login-crud-user/app/repositories/interface"
	"github.com/zidni722/login-crud-user/app/web/response"
	"gopkg.in/go-playground/validator.v9"
)

type FormLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Ctx iris.Context
	Db *gorm.DB
	Form FormLoginRequest
}

func NewLoginRequest(ctx iris.Context, db *gorm.DB, userRepository _interface.IUserRepository) LoginRequest {
	return LoginRequest{
		Ctx: ctx,
		Db: db,
	}
}

func (r *LoginRequest) Validate() bool {
	baseRequest := request.New()
	var validationErrors []string

	err := baseRequest.Validate.Struct(r.Form)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, e.Translate(baseRequest.Trans))
		}
	}

	if len(validationErrors) > 0 {
		response.ValidationResponse(r.Ctx, response.BAD_REQUEST_MESSAGE, validationErrors)
		return false
	}

	return true
}
