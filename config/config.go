package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var Cfg, _ = load()

// 配置文件
type Config struct {
	Server   `json:"server"`
	Database `json:"database"`
	Upload   `json:"upload"`
}

//服务配置
type Server struct {
	Address       string `json:"address" yaml:"address"`
	Port          string `json:"port" yaml:"port"`
	SecretKey     string `json:"secretKey" yaml:"secretKey"`
	JwtExpireTime uint   `json:"jwtExpireTime" yaml:"jwtExpireTime"`
}

// 数据库配置
type Database struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Upload struct {
	Avatar string `json:"avatar"`
	File   string `json:"file"`
}

func DB(dbCfg Database) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbCfg.User, dbCfg.Password, dbCfg.Address, dbCfg.Port, dbCfg.Name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}

func load() (cfg Config, err error) {
	var in []byte
	basePath, _ := os.Getwd()
	in, err = os.ReadFile(basePath + "/config/config.yaml")

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
