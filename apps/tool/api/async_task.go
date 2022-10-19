package api

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	"gin-bee/apps/tool/model"
	"gin-bee/apps/tool/request"
	async_task2 "gin-bee/async_task"
	"gin-bee/middleware"
	"gin-bee/utils"
	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

var CTask = TaskController{}

type TaskController struct {
}

// Create
// @Summary
// @Schemes
// @Description 执行异步任务
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.AddParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/tool/async_task [post]
func (t *TaskController) Create(c *gin.Context) {
	var param request.AddParam
	var task model.Task
	var eta *time.Time
	user, err1 := middleware.GetCurrentUser(c)
	if err1 != nil {
		c.JSONP(http.StatusUnauthorized, err1)
		return
	}
	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数不正确1"})
		return
	}
	err = c.ShouldBindBodyWith(&task, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}

	task.Creator = &user.Id
	task.State = tasks.StatePending
	task.Uid = fmt.Sprintf("task_%v", uuid.New().String())
	eta = param.Time
	addTask0 := &tasks.Signature{
		Name: "TestAdd",
		Args: []tasks.Arg{
			{
				Type:  "int64",
				Value: 1,
			},
			{
				Type:  "int64",
				Value: 2,
			},
		},
		UUID:         task.Uid,
		RetryTimeout: 3,
		RetryCount:   3,
		ETA:          eta,
	}
	task.RegisterName = addTask0.Name
	err = async_task2.AsyncTestSend(async_task2.Ser, addTask0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("创建任务失败:%s", err.Error())})
		return
	}
	// 当任务执行时间在创建之前时，任务执行时间为当前时间
	now := time.Now()
	if eta.Sub(now) < 0 {
		task.Time = &now
	}
	if err = apps.Db.Create(&task).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("创建任务失败,系统已存在任务：%s", param.Name)})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("创建任务失败:%s", err.Error())})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"msg": "创建任务成功！"})

}

// List
// @Summary
// @Schemes
// @Description 获取任务列表
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.TaskResponse
// @Failure 400 {object} response.Response
// @Router /api/tool/async_task [get]
func (t *TaskController) List(c *gin.Context) {
	var asyncTasks []model.Task
	var res []map[string]any
	if err := apps.Db.Find(&asyncTasks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "查询数据失败"})
		return
	}
	marshal, err := json.Marshal(asyncTasks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "数据解析失败"})
		return
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "数据反解析失败"})
		return
	}
	for _, item := range res {
		item["createdAt"] = utils.StrTimeFormat(item["createdAt"].(string))
		item["time"] = utils.StrTimeFormat(item["time"].(string))
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
