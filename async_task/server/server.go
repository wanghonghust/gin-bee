package server

import (
	"gin-bee/apps/async_task/task"
	"github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	"github.com/RichardKnop/machinery/v2/config"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
)

func StartServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Broker:          "redis://121.4.61.20:6379",
		ResultBackend:   "redis://121.4.61.20:6379",
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}
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
