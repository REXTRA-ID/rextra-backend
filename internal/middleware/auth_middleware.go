package middleware

import (
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
		ctx.Set("membership", idToken["membership"])
		ctx.Next()
	}
}
