package main

import (
	"fmt"
	"gin-bee/apps/system"
	"gin-bee/apps/tool"
	"gin-bee/async_task"
	_ "gin-bee/async_task"
	"gin-bee/async_task/server"
	"gin-bee/config"
	_ "gin-bee/docs"
	"gin-bee/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Bee Admin API
// @version 0.0.1
// @description Bee Admin
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	// 运行异步任务worker
	var err error
	async_task.Ser, err = server.StartServer()
	defer async_task.Ser.GetBroker().StopConsuming()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go func() {
		err := async_task.Worker(async_task.Ser)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(gin.Recovery(), middleware.CORSMiddleware(), middleware.LogMiddleware())
	system.RouterHandler(r)
	tool.RouterHandler(r)
	err = r.Run(fmt.Sprintf("%s:%s", config.Cfg.Server.Address, config.Cfg.Server.Port))
	if err != nil {
		fmt.Println(err)
	}

}
