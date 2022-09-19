package server

import (
	"gin-bee/async_task/task"
	"gin-bee/config"
	"github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
)

func StartServer() (*machinery.Server, error) {
	cnf := config.Cfg.Machinery
	broker := redisbroker.NewGR(cnf, []string{"121.4.61.20:6379"}, 0)
	backend := redisbackend.NewGR(cnf, []string{"121.4.61.20:6379"}, 0)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)
	// Register tasks
	tasks := map[string]interface{}{
		"TestAdd":        task.Add,
		"UpdateTaskInfo": task.UpdateTaskInfo,
	}
	return server, server.RegisterTasks(tasks)
}
