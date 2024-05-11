package handlers

import (
	"mim/pkg/code"
	"mim/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			ResponseError(c, code.CodeUnAuth)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) > 2 || parts[0] != "Bearer" {
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
