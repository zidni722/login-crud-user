package controllers

import (
	"github.com/jinzhu/gorm"
	golog "github.com/kataras/golog"
	"github.com/kataras/iris"
	user_request "github.com/zidni722/login-crud-user/app/dto/request/crud"
	"github.com/zidni722/login-crud-user/app/models"
	_interface "github.com/zidni722/login-crud-user/app/repositories/interface"
	"github.com/zidni722/login-crud-user/app/utils"
	"github.com/zidni722/login-crud-user/app/web/response"
)

type UserController struct {
	Db             *gorm.DB
	UserRepository _interface.IUserRepository
}

func NewUserController(db *gorm.DB, userRepository _interface.IUserRepository) *UserController {
	return &UserController{
		Db:             db,
		UserRepository: userRepository,
	}
}

func (c *UserController) CreateUserHandler(ctx iris.Context) {
	tx := c.Db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response.InternalServerErrorResponse(ctx, r)
			return
		}
	}()

	formRequest := user_request.NewUserRequest(ctx, c.Db, c.UserRepository)

	if err := ctx.ReadJSON(&formRequest.Form); err != nil {
		response.InternalServerErrorResponse(ctx, err)
		return
	}

	if !formRequest.Validate() {
		return
	}

	var user models.User

	golog.Info(formRequest)

	user.Username = formRequest.Form.Username
	user.Email = formRequest.Form.Email
	user.Address = formRequest.Form.Address
	user.Password, _ = utils.HashPassword(formRequest.Form.Password)

	if err := c.UserRepository.Create(c.Db, &user); err != nil {
		tx.Rollback()
		response.InternalServerErrorResponse(ctx, err)
		return
	}

	tx.Commit()
	response.SuccessResponse(ctx, response.CREATED, response.SUCCESS_SAVE_USER, nil)
}

func (c *UserController) GetIndexHandler(ctx iris.Context) {
	var users []models.User

	c.UserRepository.FindAll(c.Db, &users)

	if len(users) == 0 {
		response.ErrorResponse(ctx, response.UNPROCESSABLE_ENTITY, "User doesn't exists.")
		return
	}

	userResponse := response.NewUserResponse(c.Db)
	result := userResponse.Collection(users)

	response.SuccessResponse(ctx, response.OK, response.OK_MESSAGE, result)
}

func (c *UserController) GetDetailHandler(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("userId")

	var user models.User

	c.UserRepository.FindById(c.Db, &user, int(id))

	if user == (models.User{}) {
		response.ErrorResponse(ctx, response.UNPROCESSABLE_ENTITY, "User doesn't exists.")
		return
	}

	userResponse := response.NewUserResponse(c.Db)
	result := userResponse.New(user)

	response.SuccessResponse(ctx, response.OK, response.OK_MESSAGE, result)
}

func (c *UserController) UpdateUserHandler(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("userId")
	formRequest := user_request.NewUserUpdateRequest(ctx, c.Db, c.UserRepository)

	if err := ctx.ReadJSON(&formRequest.Form); err != nil {
		response.InternalServerErrorResponse(ctx, err)
		return
	}

	if !formRequest.Validate() {
		return
	}

	var user models.User

	c.UserRepository.FindById(c.Db, &user, int(id))

	if user == (models.User{}) {
		response.ErrorResponse(ctx, response.UNPROCESSABLE_ENTITY, "User doesn't exists.")
		return
	}

	if formRequest.Form.Username != "" {
		user.Username = formRequest.Form.Username
	}

	if formRequest.Form.Email != "" {
		user.Email = formRequest.Form.Email
	}

	if formRequest.Form.Address != "" {
		user.Address = formRequest.Form.Address
	}

	if !utils.CheckPasswordHash(formRequest.Form.Password, user.Password) {
		user.Password, _ = utils.HashPassword(formRequest.Form.Password)
	}

	if err := c.UserRepository.Update(c.Db, &user); err != nil {
		response.InternalServerErrorResponse(ctx, err)
		return
	}

	c.UserRepository.FindById(c.Db, &user, int(id))

	userResponse := response.NewUserResponse(c.Db)
	result := userResponse.New(user)

	response.SuccessResponse(ctx, response.OK, response.OK_MESSAGE, result)
}

func (c *UserController) DeleteUserHandler(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("userId")

	var user models.User

	c.UserRepository.FindById(c.Db, &user, int(id))

	if user == (models.User{}) {
		response.ErrorResponse(ctx, response.UNPROCESSABLE_ENTITY, "User is doesn't exists.")
		return
	}

	c.UserRepository.Delete(c.Db, &user)

	response.SuccessResponse(ctx, response.OK, response.OK_MESSAGE, nil)
}
