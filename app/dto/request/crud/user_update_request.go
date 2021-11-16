package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/zidni722/login-crud-user/app/dto/request"
	_interface "github.com/zidni722/login-crud-user/app/repositories/interface"
	"github.com/zidni722/login-crud-user/app/web/response"
	"gopkg.in/go-playground/validator.v9"
)

type FormUserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"omitempty,email"`
	Address  string `json:"address"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Ctx            iris.Context
	Db             *gorm.DB
	Form           FormUserUpdateRequest
	UserRepository _interface.IUserRepository
}

func NewUserUpdateRequest(ctx iris.Context, db *gorm.DB, userRepository _interface.IUserRepository) UserUpdateRequest {
	return UserUpdateRequest{
		Ctx:            ctx,
		Db:             db,
		UserRepository: userRepository,
	}
}

func (r *UserUpdateRequest) Validate() bool {
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
