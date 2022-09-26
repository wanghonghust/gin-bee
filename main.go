package main

import (
	_ "embed"
	"fmt"
	"gin-bee/apps"
	"gin-bee/apps/system"
	systemmodel "gin-bee/apps/system/model"
	"gin-bee/apps/tool"
	toolmodel "gin-bee/apps/tool/model"
	"gin-bee/async_task"
	"gin-bee/async_task/server"
	"gin-bee/config"
	_ "gin-bee/docs"
	"gin-bee/middleware"
	"gin-bee/utils"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli"
	"os"
)

var (
	app *cli.App
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "bee"
	app.Usage = "通用后台管理系统"
	app.Version = "1.0.0"
}

// @title Bee Admin API
// @version 0.0.1
// @description Bee Admin
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	setCommand()
}
func startServer() (err error) {
	// 运行异步任务worker
	async_task.Ser, err = server.StartServer()
	defer async_task.Ser.GetBroker().StopConsuming()
	if err != nil {
		zaplog.Logger.Error(err)
		return
	}
	go func() {
		err := async_task.Worker(async_task.Ser)
		if err != nil {
			zaplog.Logger.Error(err)
			return
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
		zaplog.Logger.Error(err)
		return
	}
	return nil
}

func initDataBase() (err error) {
	// 初始化数据库
	err = systemmodel.InitFileMenu()
	if err != nil {
		return err
	}
	err = systemmodel.InitUser()
	if err != nil {
		return err
	}
	err = systemmodel.InitRolePerM()
	if err != nil {
		return err
	}
	err = systemmodel.InitAPI()
	if err != nil {
		return err
	}
	err = systemmodel.InitLog()
	if err != nil {
		return err
	}
	err = toolmodel.InitTask()
	if err != nil {
		return err
	}
	return nil
}

func setCommand() {
	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "init database",
			Action: func(c *cli.Context) error {
				if err := initDataBase(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "server",
			Usage: "start bee server ",
			Action: func(c *cli.Context) error {
				if err := startServer(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "createuser",
			Usage: "create a normal account ",
			Action: func(c *cli.Context) (err error) {
				err = createUser(false)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				fmt.Println("Success！")
				return nil
			},
		},
		{
			Name:  "createsuperuser",
			Usage: "create a super user account ",
			Action: func(c *cli.Context) (err error) {
				err = createUser(true)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				fmt.Println("Success！")
				return nil
			},
		}, {
			Name:  "createrole",
			Usage: "create an role",
			Action: func(c *cli.Context) (err error) {
				err = createRole("admin")
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				fmt.Println("Success！")
				return nil
			},
		},
	}

	// Run the CLI app
	_ = app.Run(os.Args)
}

func createUser(isSuperUser bool) (err error) {
	var username string
	var password string
	fmt.Print("Username:")
	_, err = fmt.Scanln(&username)
	if err != nil {
		return err
	}
	fmt.Print("Password:")
	masked, err := gopass.GetPasswdMasked()
	if err != nil {
		return err
	}
	password = string(masked)
	hashPwd, err := utils.Password(password)
	if err != nil {
		return err
	}
	user := systemmodel.User{Username: username, Password: hashPwd, IsSuperUser: isSuperUser}
	if err = apps.Db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func createRole(name string) (err error) {
	role := systemmodel.Role{Name: name}
	if err = apps.Db.Create(&role).Error; err != nil {
		return err
	}
	fmt.Println(role)
	return nil
}
