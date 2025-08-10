package config

import (
	"fmt"
	"rextra-backend/internal/middleware"
	mylog "rextra-backend/internal/pkg/logger"
	"rextra-backend/internal/pkg/response"

	"net/http"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func NewRouter(server *gin.Engine) *gin.Engine {
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
		})
	})

	server.MaxMultipartMemory = 30 * 1024 * 1024
	server.Use(customRecovery())
	server.Use(middleware.CORSMiddleware())

	server.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 123",
		})
	})

	return server
}

func customRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var wrappedErr error
				if e, ok := err.(error); ok {
					wrappedErr = e
				} else {
					wrappedErr = fmt.Errorf("%v", err)
				}

				fmt.Println(mylog.ColorizePanic(fmt.Sprintf("\n[Recovery] Panic occurred: %v\n", err)))
				stack := debug.Stack()
				coloredStack := mylog.ColorizePanic(string(stack))

				fmt.Fprintln(os.Stderr, coloredStack)
				response.NewFailed("server panic occured", wrappedErr).
					SendWithAbort(ctx)
			}
		}()

		ctx.Next()
	}
}
