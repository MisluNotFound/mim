package main

import (
	"fmt"
	"mim/db"
	"mim/internal/api"
	"mim/internal/connect"
	"mim/internal/logic"
	"mim/pkg/logger"
	"mim/pkg/snowflake"
	"mim/setting"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	setting.Init("./conf/config.yaml")

	logger.Init(setting.Conf.LogConfig, setting.Conf.Mode)
	snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID)
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", setting.Conf.MySQLConfig.User, setting.Conf.MySQLConfig.Password, setting.Conf.MySQLConfig.Host, setting.Conf.MySQLConfig.Port, setting.Conf.MySQLConfig.DB)
	db.InitDB(dsn)
	db.InitRDB(setting.Conf.RedisConfig.Addr)
	
	logic.InitLogic()
	connect.InitConnect()
	api.InitAPI()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
