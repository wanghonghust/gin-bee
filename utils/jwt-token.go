package utils

import (
	"gin-bee/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type UserInfo struct {
	Id       uint
	UserName string
	State    bool
}

type JwtClaims struct {
	UserInfo
	jwt.RegisteredClaims
}

var (
	tokenExpireDuration = time.Second * time.Duration(config.Cfg.Server.JwtExpireTime)
	secretKey           = config.Cfg.Server.SecretKey
)

// GenerateToken 生成token字符串
func GenerateToken(userInfo UserInfo) (string, error) {
	expirationTime := time.Now().Add(tokenExpireDuration)
	claims := &JwtClaims{
		userInfo,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			Issuer:    "sever",
		},
	}
	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := Token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解析token字符串
func ParseToken(token string) (*JwtClaims, error) {
	claim := &JwtClaims{}
	_, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil

	})
	if err != nil {
		return nil, err
	}
	return claim, nil
}
