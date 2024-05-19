// api层负责为客户端提供基本服务的接口并交由logic层处理返回响应结果
package api

import (
	"mim/internal/api/handlers"
	"mim/internal/api/rpc"
	"mim/pkg/logger"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitAPI() {
	go rpc.InitAPIRpc()
	go router()
}

func router() {
	r = gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	u := r.Group("/user")
	{
		u.POST("/signin", handlers.SignIn)
		u.POST("/signup", handlers.SignUp)
	}

	g := r.Group("/group").Use(handlers.Auth())
	{
		g.POST("/new", handlers.NewGroup)
		g.POST("/join", handlers.JoinGroup)
		g.GET("/get", handlers.FindGroup)
		g.DELETE("/leave", handlers.LeaveGroup)
	}

	m := r.Group("/message").Use(handlers.Auth())
	{
		m.GET("/pull", handlers.PullMessage)
	}
	r.Run(":8080")
}
