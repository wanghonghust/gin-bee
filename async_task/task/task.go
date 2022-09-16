package task

import (
	"fmt"
)

func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

func UpdateTaskInfo(taskUid string, taskRes ...any) error {
	fmt.Println("这是在更新任务的信息：", taskUid, taskRes)
	return nil
}
