package middleware

import (
	"github.com/zidni722/login-crud-user/app/utils"
	"github.com/zidni722/login-crud-user/app/web/response"
	"github.com/kataras/iris"
	"strings"
)

type AuthMiddleware struct {
	Authable utils.Authable
}

func NewAuthMiddleware(authable utils.Authable) *AuthMiddleware {
	return &AuthMiddleware{
		Authable: authable,
	}
}

func (m *AuthMiddleware) AuthRequired(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		arrayAuthHeader := strings.Split(authHeader, " ")
		if len(arrayAuthHeader) == 2 {
			if arrayAuthHeader[0] == "Bearer" {
				tokenString := arrayAuthHeader[1]
				user, err := m.Authable.Decode(tokenString)

				if (user != nil) && err == nil {
					ctx.Values().Set("user", *user)
					ctx.Next()
					return
				}
			}
		}
	}

	response.UnAuthorizedResponse(ctx)
	ctx.StopExecution()
	return
}
