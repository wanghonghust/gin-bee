package auth

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	"gin-bee/utils"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

var cAuth = Auth{}

type Auth struct {
}

type User struct {
	Id       int
	Username string
	Password string
}

type Users struct {
	Id []int
}

type Userinfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username" `
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Avatar    uint      `json:"avatar"`
	State     bool      `json:"state"`
}

type PwdChange struct {
	Id           int    `json:"id"`
	OldPassword  string `json:"oldPassword"`
	NewPassword  string `json:"newPassword"`
	ConfPassword string `json:"confPassword"`
}

func init() {
	defer func() {
		zaplog.Logger.Info("数据表user迁移成功")
	}()
	user := UserModel{}
	err := user.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}
}
func (a *Auth) Auth(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		c.JSONP(http.StatusUnauthorized, gin.H{
			"code": 1,
			"msg":  "未登录状态！",
		})
		return
	}
	if strings.HasPrefix(authorization, "Bearer ") {
		resJwt := strings.Split(authorization, " ")
		token := resJwt[len(resJwt)-1]
		_, err := utils.ParseToken(token)
		if err != nil {
			fmt.Println(err)
			c.JSONP(http.StatusUnauthorized, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			return
		}
	} else {
		c.JSONP(http.StatusUnauthorized, gin.H{
			"code": 1,
			"msg":  "token格式不正确",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"auth": true})
}

func (a *Auth) AllUser(c *gin.Context) {
	// 包含分页
	paginator := utils.Paginator{}
	var users []Userinfo
	var usersData []interface{}
	apps.Db.Model(UserModel{}).Find(&users)
	err := c.BindJSON(&paginator)
	marshal, err := json.Marshal(users)
	err = json.Unmarshal(marshal, &usersData)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	err = paginator.Init(usersData, paginator.PerPage)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	res, err := paginator.PageData(paginator.Page)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, res)

}

func (a *Auth) Login(c *gin.Context) {
	// 用户登录
	user := User{}
	userModel := UserModel{}
	err := c.BindJSON(&user)
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}

	tx := userModel.Where("username = ?", user.Username).First(&userModel)
	if tx.Error == gorm.ErrRecordNotFound {
		c.JSONP(http.StatusUnauthorized, gin.H{"code": 1, "msg": "账户不存在"})
		return
	}
	if userModel.State == false {
		c.JSONP(http.StatusForbidden, gin.H{"code": 1, "msg": "用户被禁用"})
		return
	}
	if utils.PasswordVerify(user.Password, userModel.Password) {
		// 保存session
		userInfo := utils.UserInfo{Id: userModel.ID, UserName: userModel.Username, State: userModel.State}
		token, err := utils.GenerateToken(userInfo)
		if err != nil {
			zaplog.Logger.Error(err.Error())
		}
		c.JSONP(http.StatusOK, gin.H{"code": 0, "data": gin.H{
			"user":     userModel.ID,
			"username": userModel.Username,
			"avatar":   userModel.Avatar,
			"email":    userModel.Email,
			"createAt": userModel.CreatedAt,
			"token":    token,
		}})
		return
	} else {
		c.JSONP(http.StatusUnauthorized, gin.H{"code": 1, "msg": "密码错误"})
		return
	}

}

func (a *Auth) CreateUser(c *gin.Context) {
	// 用户注册
	var userMap map[string]any
	queryUser := UserModel{}

	err := c.ShouldBindBodyWith(&queryUser, binding.JSON)
	err = c.ShouldBindBodyWith(&userMap, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}

	if userMap["password"] == userMap["password_c"] {
		if queryUser.Password == "" {
			c.JSONP(http.StatusOK, gin.H{"code": 0, "msg": "密码为空"})
			return
		}
		// 密码加密
		encryptPwd, _ := utils.Password(queryUser.Password)
		queryUser.Password = encryptPwd
		tx := apps.Db.Create(&queryUser)
		if tx.Error != nil {
			if strings.Contains(tx.Error.Error(), "Duplicate entry") {
				c.JSONP(http.StatusOK, gin.H{"code": 0, "msg": "已有相同的用户名"})
				return
			} else {
				c.JSONP(http.StatusOK, gin.H{"code": 0, "msg": "未知错误"})
				return
			}
		}
		c.JSONP(http.StatusOK, gin.H{"code": 1, "msg": "创建用户成功"})
		return
	} else {
		c.JSONP(http.StatusOK, gin.H{"code": 0, "msg": "两次输入密码不相同"})
		return
	}
}

func (a *Auth) RegisterIndex(c *gin.Context) {
	// 用户注册页面
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "注册"})
}

func (a *Auth) loginIndex(c *gin.Context) {
	// 用户登录页面
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "登录"})
}

func (a *Auth) UserInfo(c *gin.Context) {
	// 用户信息
	userModel := UserModel{}
	err := c.BindJSON(&userModel)
	if err != nil {
		fmt.Println(err)
	}
	userModel.Where("id = ?", userModel.ID).First(&userModel)
	byteinfo, _ := json.Marshal(&userModel)
	userinfo := make(gin.H)
	err = json.Unmarshal(byteinfo, &userinfo)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	userinfo["created_at"] = utils.StrTimeFormat(utils.String(userinfo["created_at"]))
	delete(userinfo, "file")
	delete(userinfo, "password")
	c.JSONP(http.StatusOK, userinfo)
}
func (a *Auth) UpdateUserInfo(c *gin.Context) {
	queryUser := UserModel{}
	var bodyMap map[string]any
	user := UserModel{}
	err := c.ShouldBindBodyWith(&queryUser, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	err = c.ShouldBindBodyWith(&bodyMap, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	tx := apps.Db.First(&user, queryUser.ID)
	if tx.Error != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}

	if bodyMap["state"] != nil {
		tx = tx.Select("state").Updates(queryUser)
	} else {
		tx = tx.Select("nickname", "email").Updates(queryUser)
	}

	if tx.Error != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "更新失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "请求成功"})
}
func (a *Auth) EditUserAvatar(c *gin.Context) {
	// 编辑用户信息
	var err error
	var user UserModel

	err = c.BindJSON(&user)
	if err != nil || user.ID == 0 || user.Avatar == nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	tx := apps.Db.Model(&user).Where("id = ?", user.ID).Update("avatar", user.Avatar)
	if tx.Error != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "修改失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"code": 0, "msg": "修改成功"})
}

func (a *Auth) DeleteUSer(c *gin.Context) {
	var users Users
	err := c.BindJSON(&users)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	tx := apps.Db.Where("id IN (?)", users.Id).Delete(&UserModel{})
	if tx.Error != nil {
		c.JSONP(http.StatusExpectationFailed, gin.H{"msg": "删除失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "删除成功"})
}

func (a *Auth) ChangePwd(c *gin.Context) {
	var user PwdChange
	var err error
	err = c.BindJSON(&user)
	if err != nil {
		fmt.Println(err)
	}
	userModel, msg := user.Validate()
	hashPwd, _ := utils.Password(user.NewPassword)
	if msg != "" {
		c.JSONP(http.StatusBadRequest, gin.H{"code": 1, "msg": msg})
		return
	}
	apps.Db.Model(&userModel).Where("id = ?", user.Id).Update("password", hashPwd)
	c.JSONP(http.StatusOK, gin.H{"code": 0, "msg": "修改成功"})
}

func (c *PwdChange) Validate() (user *UserModel, msg string) {
	msg = ""

	tx := apps.Db.Model(&user).Where("id = ?", c.Id).First(&user)
	if tx.Error == gorm.ErrRecordNotFound {
		msg = "用户未找到"
		return
	}
	if !utils.PasswordVerify(c.OldPassword, user.Password) {
		msg = "密码不正确"
		return
	}
	if len(c.NewPassword) < 8 || len(c.NewPassword) > 20 {
		msg = "密码长度应大于等于8位小于等于20位"
		return
	}
	if c.NewPassword != c.ConfPassword {
		msg = "确认密码等于新密码"
		return
	}
	return
}
