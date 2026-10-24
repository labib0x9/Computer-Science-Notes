package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		slog.Info("Incoming request",
			"Method", ctx.Request.Method,
			"Path", ctx.Request.URL.Path,
			"Remote-addr", ctx.Request.RemoteAddr,
		)
	}
}

func pingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// For request
		// ctx e thakbe
		ctx.Set("X-XXX", "i dont know")

		// Header e jbe
		ctx.Header("X-Hackerone-Id", "9xzer0")


		// For response
		ctx.Writer.Header().Set("Resp-Header", "labib")

		ctx.Next()
	}
}

func pingHandler(ctx *gin.Context) {
	msg := map[string]string{
		"msg":  "pong",
		"ping": "pingpong",
	}
	ctx.JSON(http.StatusOK, msg)
}

func main() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// middleware
	router.Use(Logger())
	router.Use(pingMiddleware())

	// path
	router.GET("/ping", pingHandler)

	slog.Info("Started Server at 127.0.0.1:8080")
	router.Run(":8080")
}
