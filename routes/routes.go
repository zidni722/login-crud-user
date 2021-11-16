package routes

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/zidni722/login-crud-user/app/repositories/impl"
	"github.com/zidni722/login-crud-user/app/utils"
	"github.com/zidni722/login-crud-user/app/web/controllers"
	"github.com/zidni722/login-crud-user/app/web/middleware"
	"github.com/zidni722/login-crud-user/bootstrap"
	"github.com/zidni722/login-crud-user/config"
)

type Route struct {
	Config      *config.Configuration
	CorsHandler context.Handler
}

func NewRoute(config *config.Configuration) *Route {
	return &Route{
		Config: config,
	}
}

func (r *Route) Configure(b *bootstrap.Bootstrapper) {
	b.Get("/", controllers.GetHomeHandler)

	userRepository := impl.NewUserRepositoryImpl()

	jwtService := utils.NewJWTService(r.Config.Database.DB, userRepository)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	user := b.Party("/users").AllowMethods(iris.MethodOptions)
	{
		userController := controllers.NewUserController(r.Config.Database.DB, userRepository)
		user.Use(authMiddleware.AuthRequired)
		{
			user.Get("/", userController.GetIndexHandler)
			user.Post("/", userController.CreateUserHandler)
			user.Put("/{userId:uint}", userController.UpdateUserHandler)
			user.Delete("/{userId:uint}", userController.DeleteUserHandler)
			user.Get("/{userId:uint}", userController.GetDetailHandler)
		}
	}

	loginController := controllers.NewLoginController(r.Config.Database.DB, userRepository, jwtService)
	b.Post("/login", loginController.LoginHandler)
}
