package middleware

import (
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// OnlyAllow permits access when the user's role matches one of the allowed roles.
func (m Middleware) OnlyAllow(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleValue, exists := ctx.Get("role")
		role, ok := roleValue.(string)
		if !exists || !ok {
			response.NewFailed("forbidden", myerror.RoleNotAllowed()).SendWithAbort(ctx)
			return
		}

		for _, allowed := range roles {
			if role == allowed {
				ctx.Next()
				return
			}
		}

		response.NewFailed("forbidden", myerror.RoleNotAllowed()).SendWithAbort(ctx)
	}
}

// OnlyAdmin ensures only ADMIN role can access the route.
// Equivalent to OnlyAllow("ADMIN").
func (m Middleware) OnlyAdmin() gin.HandlerFunc {
	return m.OnlyAllow("ADMIN")
}

// OnlyUser ensures only USER role can access the route.
// Equivalent to OnlyAllow("USER").
func (m Middleware) OnlyUser() gin.HandlerFunc {
	return m.OnlyAllow("USER")
}
