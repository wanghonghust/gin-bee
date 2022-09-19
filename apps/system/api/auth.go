package api

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	system "gin-bee/apps/system/model"
	"gin-bee/apps/system/request"
	"gin-bee/middleware"
	"gin-bee/response"
	"gin-bee/utils"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var CAuth = Auth{}

type Auth struct {
}

type User struct {
	Id       int
	Username string
	Password string
}

type Users struct {
	Id []uint
}

type Userinfo struct {
	ID        uint                `json:"id"`
	CreatedAt time.Time           `json:"createdAt"`
	Username  string              `json:"username" `
	Nickname  string              `json:"nickname"`
	Email     string              `json:"email"`
	Avatar    uint                `json:"avatar"`
	State     bool                `json:"state"`
	Menu      []response.TreeMenu `json:"menu"`
}

type PwdChange struct {
	Id           int    `json:"id"`
	OldPassword  string `json:"oldPassword"`
	NewPassword  string `json:"newPassword"`
	ConfPassword string `json:"confPassword"`
}

// Auth
// @Summary
// @Schemes
// @Description 登录验证
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth [post]
func (a *Auth) Auth(c *gin.Context) {
	_, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSONP(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"auth": true})
}

// AllUser
// @Summary
// @Schemes
// @Description 获取所有用户
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/user [post]
func (a *Auth) AllUser(c *gin.Context) {
	// 包含分页
	paginator := utils.Paginator{}
	var count int64
	var users []system.User
	var usersData []map[string]any
	err := c.ShouldBindBodyWith(&paginator, binding.JSON)
	if paginator.Page <= 0 {
		apps.Db.Model(system.User{}).Preload("Role").Find(&users)
	} else {
		apps.Db.Model(system.User{}).Preload("Role").Omit("password").Offset((paginator.Page - 1) * paginator.PerPage).Limit(paginator.PerPage).Find(&users)
	}
	tx := apps.Db.Model(system.User{}).Count(&count)
	if tx.Error != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询错误"})
		return
	}

	marshal, err := json.Marshal(users)
	err = json.Unmarshal(marshal, &usersData)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	for idx, _ := range usersData {
		usersData[idx]["createdAt"] = utils.StrTimeFormat(fmt.Sprintf("%v", usersData[idx]["createdAt"]))
	}
	err = paginator.Init(usersData, paginator.PerPage, uint(count))
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

// Login
// @Summary
// @Schemes
// @Description 登录
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/login [post]
func (a *Auth) Login(c *gin.Context) {
	// 用户登录
	user := User{}
	userModel := system.User{}
	err := c.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}

	tx := userModel.Where("username = ?", user.Username).First(&userModel)
	if tx.Error == gorm.ErrRecordNotFound {
		c.JSONP(http.StatusUnauthorized, gin.H{"code": 01, "msg": "账户不存在"})
		return
	}
	if userModel.State == false {
		c.JSONP(http.StatusForbidden, gin.H{"code": 02, "msg": "用户被禁用"})
		return
	}
	if utils.PasswordVerify(user.Password, userModel.Password) {
		// 保存session
		userInfo := utils.UserInfo{Id: userModel.ID, UserName: userModel.Username, State: userModel.State}
		token, err := utils.GenerateToken(userInfo)
		if err != nil {
			zaplog.Logger.Error(err.Error())
		}
		c.JSONP(http.StatusOK, gin.H{"code": 00, "data": gin.H{
			"user":     userModel.ID,
			"username": userModel.Username,
			"avatar":   userModel.Avatar,
			"email":    userModel.Email,
			"createAt": userModel.CreatedAt,
			"token":    token,
		}})
		return
	} else {
		c.JSONP(http.StatusUnauthorized, gin.H{"code": 03, "msg": "密码错误"})
		return
	}

}

// CreateUser
// @Summary
// @Schemes
// @Description 创建用户
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/user/create [post]
func (a *Auth) CreateUser(c *gin.Context) {
	// 用户注册
	var param request.CreateUserParam
	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	queryUser := system.User{
		Username: param.Username,
		Password: param.Password,
		Nickname: param.Nickname,
		Email:    param.Email,
		Role: func(nums []uint) (res []system.Role) {
			for _, item := range nums {
				res = append(res, system.Role{Model: apps.Model{
					ID: item,
				}})
			}
			return
		}(param.Role),
	}
	if param.Password == param.PasswordC {
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

// UserInfo
// @Summary
// @Schemes
// @Description 用户信息
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/user [get]
func (a *Auth) UserInfo(c *gin.Context) {
	// 用户信息
	var menus []system.Menu
	id := c.Query("id")
	uid, _ := strconv.Atoi(id)
	userModel := system.User{Model: apps.Model{ID: uint(uid)}}
	userModel.First(&userModel)
	byteinfo, _ := json.Marshal(&userModel)
	userinfo := make(gin.H)
	err := json.Unmarshal(byteinfo, &userinfo)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}

	apps.Db.Joins("left join role_menus on menus.id = role_menus.menu_id ").
		Joins("left join user_roles on role_menus.role_id = user_roles.role_id").
		Joins("left join users on user_roles.user_id = users.id").Where("users.id = ?", userModel.ID).Find(&menus)

	treeMenu, err := TreeOfMenus(menus)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询错误"})
		return
	}

	userinfo["createdAt"] = utils.StrTimeFormat(utils.String(userinfo["createdAt"]))
	userinfo["menu"] = treeMenu
	delete(userinfo, "file")
	delete(userinfo, "password")
	c.JSONP(http.StatusOK, userinfo)
}

// UpdateUserInfo
// @Summary
// @Schemes
// @Description 编辑用户
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/user [put]
func (a *Auth) UpdateUserInfo(c *gin.Context) {
	queryUser := request.UpdateUserParam{}
	var bodyMap map[string]any
	user := system.User{}
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
	user.ID = queryUser.ID
	user.Email = queryUser.Email
	user.Nickname = queryUser.Nickname
	user.Role = func(nums []uint) (res []system.Role) {
		for _, item := range nums {
			res = append(res, system.Role{Model: apps.Model{
				ID: item,
			}})
		}
		return
	}(queryUser.Role)

	if _, ok := bodyMap["state"]; ok {
		err = apps.Db.Select("state").Updates(&user).Error
		if err != nil {
			c.JSONP(http.StatusBadRequest, gin.H{"msg": "更新失败"})
			return
		}
	} else {
		apps.Db.Select("nickname", "email").Updates(&user)
		err = apps.Db.Model(&user).Association("Role").Replace(user.Role)
		if err != nil {
			c.JSONP(http.StatusBadRequest, gin.H{"msg": "更新失败"})
			return
		}
	}
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "更新失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "请求成功"})
}

// EditUserAvatar
// @Summary
// @Schemes
// @Description 编辑用户头像
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/user/avatar [put]
func (a *Auth) EditUserAvatar(c *gin.Context) {
	// 编辑用户信息
	var err error
	var user system.User

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

// DeleteUSer
// @Summary
// @Schemes
// @Description 删除用户
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id body int true "用户id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/user [delete]
func (a *Auth) DeleteUSer(c *gin.Context) {
	var param Users
	err := c.BindJSON(&param)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
		return
	}
	users := make([]system.User, len(param.Id))
	for _, item := range param.Id {
		users = append(users, system.User{Model: apps.Model{ID: item}})
	}

	// 开启事务
	err = apps.Db.Transaction(func(tx *gorm.DB) error {
		// 删除，彻底删除,并清除关联
		if err2 := tx.Unscoped().Select(clause.Associations).Delete(&users).Error; err2 != nil {
			return err2
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		c.JSONP(http.StatusExpectationFailed, gin.H{"msg": "删除失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "删除成功"})
}

// ChangePwd
// @Summary
// @Schemes
// @Description 修改用户密码
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/passwd [put]
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

func (c *PwdChange) Validate() (user system.User, msg string) {
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
