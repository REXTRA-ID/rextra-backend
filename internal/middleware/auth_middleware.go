package middleware

import (
	"fmt"
	myerror "rextra-backend/internal/pkg/error"
	myjwt "rextra-backend/internal/pkg/jwt"
	"rextra-backend/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	MESSAGE_FAILED_VERIFY_TOKEN = "failed to verify token"
	MESSAGE_USER_NOT_AUTHORIZED = "user not authorized"
	MESSAGE_API_IS_LOCKED       = "api is now locked"
)

func (m Middleware) OnlyAllow(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.MustGet("role").(string)

		for _, role := range roles {
			if userRole == role {
				ctx.Next()
				return
			}
		}

		fmt.Println(userRole)
		fmt.Println(roles)

		res := response.NewFailed(MESSAGE_USER_NOT_AUTHORIZED, myerror.RoleNotAllowed())
		res.SendWithAbort(ctx)
	}
}

func (m Middleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			res := response.NewFailed(MESSAGE_FAILED_VERIFY_TOKEN, myerror.InvalidToken())
			res.SendWithAbort(ctx)
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			res := response.NewFailed(MESSAGE_FAILED_VERIFY_TOKEN, myerror.InvalidToken())
			res.SendWithAbort(ctx)
			return
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)

		idToken, err := myjwt.GetPayloadInsideToken(authHeader)
		if err != nil {
			if err.Error() == "token expired" {
				res := response.NewFailed(MESSAGE_FAILED_VERIFY_TOKEN, myerror.InvalidToken())
				res.SendWithAbort(ctx)
				return
			}

			res := response.NewFailed(MESSAGE_FAILED_VERIFY_TOKEN, myerror.ErrGeneral)
			res.SendWithAbort(ctx)
			return
		}

		ctx.Set("token", authHeader)
		ctx.Set("payload", idToken)
		ctx.Set("user_id", idToken["user_id"])
		ctx.Set("email", idToken["email"])
		ctx.Set("role", idToken["role"])
		fmt.Println(idToken)
		ctx.Next()
	}
}
