package main

import (
	"fmt"
	"gin-bee/apps/auth"
	"gin-bee/apps/system"
	"gin-bee/config"
	"gin-bee/middleware"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {
	f, _ := os.Create("./zaplog/gin.zaplog")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	auth.RouterHandler(r)
	system.RouterHandler(r)
	err := r.Run(config.Cfg.Server.Address + ":" + config.Cfg.Server.Port)
	if err != nil {
		fmt.Println(err)
	}

}
