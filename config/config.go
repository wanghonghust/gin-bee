package config

import (
	"embed"
	"fmt"

	machineryCfg "github.com/RichardKnop/machinery/v2/config"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed config.yaml
var f embed.FS
var Cfg, _ = load()

// Config
// @Description: 配置
type Config struct {
	Server    *Server              `json:"server"`
	Database  *Database            `json:"database"`
	Upload    *Upload              `json:"upload"`
	Machinery *machineryCfg.Config `json:"machinery"`
	Redis     *redis.Options       `json:"redis"`
}

// Server
// @Description: 服务配置
type Server struct {
	Address       string `json:"address" yaml:"address"`
	Port          string `json:"port" yaml:"port"`
	SecretKey     string `json:"secretKey" yaml:"secretKey"`
	JwtExpireTime uint   `json:"jwtExpireTime" yaml:"jwtExpireTime"`
}

// Database
// @Description: 数据库配置
type Database struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Upload
// @Description: 文件上传配置
type Upload struct {
	Avatar string `json:"avatar"`
	File   string `json:"file"`
}

func DB(dbCfg Database) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&parseTime=true", dbCfg.User, dbCfg.Password, dbCfg.Address, dbCfg.Port, dbCfg.Name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db, err = gorm.Open(sqlite.Open("dsn.db"), &gorm.Config{})
	return
}

func load() (cfg Config, err error) {
	var in []byte
	// basePath, _ := os.Getwd()
	// in, err = f.ReadFile(basePath + "/config/config.yaml")
	in, err = f.ReadFile("config.yaml")

	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(in, &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
