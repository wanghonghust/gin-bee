package async_task

import (
	"encoding/json"
	"gin-bee/apps"
	"gin-bee/apps/tool/model"
	"gin-bee/zaplog"
	"github.com/RichardKnop/machinery/example/tracers"
	"github.com/RichardKnop/machinery/v2"
	backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
	"github.com/RichardKnop/machinery/v2/tasks"
)

var (
	Ser *machinery.Server
)

func Worker(server *machinery.Server) error {
	// 获取执行的任务
	var asyncTask *tasks.Signature
	consumerTag := "machinery_worker"

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		zaplog.Logger.Fatal("Unable to instantiate a tracer:" + err.Error())
	}
	defer cleanup()

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)
	backend := server.GetBackend()
	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorhandler := func(err error) {
		zaplog.Logger.Errorf("任务:%s，执行失败，ERROR:%s", asyncTask.UUID, err.Error())
		// 更新任务信息
		err1 := updateTaskInfo(backend, asyncTask)
		if err1 != nil {
			return
		}
	}

	pretaskhandler := func(signature *tasks.Signature) {
		// 记录执行的任务
		asyncTask = signature
		//zaplog.Logger.Infof("I am a start of task handler for:"+signature.Name, signature.UUID)
	}

	posttaskhandler := func(signature *tasks.Signature) {
		//zaplog.Logger.Info("I am an end of task handler for:"+signature.Name, signature.UUID)
		// 更新任务信息
		err3 := updateTaskInfo(backend, signature)
		if err3 != nil {
			return
		}

	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)
	return worker.Launch()
}

func updateTaskInfo(backend backendsiface.Backend, signature *tasks.Signature) (err error) {
	res, err := backend.GetState(signature.UUID)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(res.Results)
	if err != nil {
		return err
	}
	var task = model.Task{Uid: signature.UUID, Result: string(marshal), State: res.State}

	if err = apps.Db.Select("state", "result").Where("uid = ?", signature.UUID).Updates(&task).Error; err != nil {
		zaplog.Logger.Errorf("更新任务信息失败：%v", err)
		return err
	}
	return nil
}
