package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromCtx(ctx *gin.Context) (string, error) {
	user, exists := ctx.Get("user_id")
	if !exists {
		return "", errors.New("user id not found")
	}

	userId, ok := user.(string)
	if !ok {
		return "", errors.New("invalid user id")
	}

	return userId, nil
}
