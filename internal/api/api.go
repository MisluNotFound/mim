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
	r.Use(handlers.CorsMiddleware())
	r.GET("/upload/token", handlers.Auth(), handlers.GetOssCredentials)
	u := r.Group("/user")
	{
		u.POST("/signin", handlers.SignIn)
		u.POST("/signup", handlers.SignUp)
		u.GET("/getinfo", handlers.Auth(), handlers.GetInfo)
		u.POST("/update/password", handlers.Auth(), handlers.UpdatePassword)
		u.POST("/update/name", handlers.Auth(), handlers.UpdateName)
		u.POST("/update/photo", handlers.Auth(), handlers.UpdatePhoto)
	}

	n := r.Group("/nearby").Use(handlers.Auth())
	{
		n.POST("/open", handlers.NearbyOpen)
	}

	f := r.Group("/friend").Use(handlers.Auth())
	{
		f.POST("/add", handlers.AddFriend)
		f.GET("/get", handlers.GetFriends)
		f.DELETE("/remove", handlers.RemoveFriend)
		f.POST("/update/remark", handlers.UpdateFriendRemark)
		f.GET("/find", handlers.FindFriend)
	}

	g := r.Group("/group").Use(handlers.Auth())
	{
		g.POST("/new", handlers.NewGroup)
		g.POST("/join", handlers.JoinGroup)
		g.GET("/find", handlers.FindGroup)
		g.GET("/getall", handlers.GetGroups)
		g.POST("/leave", handlers.LeaveGroup)
		g.GET("/members", handlers.GetMembers)
		g.GET("/role", handlers.GetRole)
		g.POST("/update/photo", handlers.UpdateGroupPhoto)
		g.POST("/remove/member", handlers.RemoveMember)
	}

	m := r.Group("/message").Use(handlers.Auth())
	{
		m.GET("/pull", handlers.PullMessage)
		m.GET("/pulloffline/count", handlers.GetUnReadCount)
		m.POST("/pulloffline", handlers.PullOfflineMessage)
		m.GET("/pullerr", handlers.PullErrMessage)
	}

	r.Run(":3000")
}
