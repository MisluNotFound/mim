package handlers

import (
	"fmt"
	"mim/pkg/code"
	"mim/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println(authHeader)

		if authHeader == "" {
			ResponseError(c, code.CodeUnAuth)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) < 2 || parts[0] != "Bearer" {
			ResponseError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			ResponseError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}

		c.Set("userId", mc.UserID)
		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // 添加 Authorization 到允许的请求头部

		// 如果请求方法是 OPTIONS，则表示预检请求，直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
