package async_task

import (
	"fmt"
	"gin-bee/zaplog"
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/tasks"
)

func TestSend(server *machinery.Server, task *tasks.Signature) error {
	err := server.RegisterPeriodicTask("*/5 * * * * *", "period_task", task)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//asyncResult, _ := server.SendTask(&task)
	//fmt.Println(asyncResult)

	return nil
}

func AsyncTestSend(server *machinery.Server, task *tasks.Signature) error {
	asyncResult, _ := server.SendTask(task)
	fmt.Println(asyncResult)
	return nil
}

func AsyncChain(server *machinery.Server, task *tasks.Signature) error {
	fmt.Println("执行时间：", task.ETA)
	signature1 := tasks.Signature{
		Name: "UpdateTaskInfo",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: task.UUID,
			},
		},
	}
	chain, _ := tasks.NewChain(task, &signature1)
	res, err := server.SendChain(chain)
	if err != nil {
		zaplog.Logger.Error(err)
		return err
	}
	zaplog.Logger.Info(res)
	return nil
}
