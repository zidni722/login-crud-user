package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	user_request "github.com/zidni722/login-crud-user/app/dto/request/crud"
	"github.com/zidni722/login-crud-user/app/models"
	_interface "github.com/zidni722/login-crud-user/app/repositories/interface"
	"github.com/zidni722/login-crud-user/app/utils"
	"github.com/zidni722/login-crud-user/app/web/response"
)


type LoginController struct {
	Db *gorm.DB
	UserRepository _interface.IUserRepository
	Authable utils.Authable
}

func NewLoginController(db *gorm.DB, UserRepository _interface.IUserRepository, authable utils.Authable) *LoginController {
	return &LoginController{
		Db: db,
		UserRepository: UserRepository,
		Authable: authable,
	}
}

func (c *LoginController) LoginHandler(ctx iris.Context) {
	formRequest := user_request.NewLoginRequest(ctx, c.Db, c.UserRepository)

	if err := ctx.ReadJSON(&formRequest.Form); err != nil {
		response.InternalServerErrorResponse(ctx, err)
		return
	}

	if !formRequest.Validate() {
		return
	}

	var token string
	var err error

	var user models.User
	c.UserRepository.FindByUsername(c.Db, &user, formRequest.Form.Username)
	
	if user == (models.User{}) {
		response.UnAuthorizedResponse(ctx)
		return
	}

	token, err = c.Authable.Encode(int(user.ID), *&user.Username)
	
	if err != nil {
		response.InternalServerErrorResponse(ctx, err)
		return
	}

	loginResponse := response.NewLoginResponse(c.Db)
	result := loginResponse.New(token)

	response.SuccessResponse(ctx, response.OK, response.OK_MESSAGE, result)
}
