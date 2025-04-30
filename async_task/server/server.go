package server

import (
	"fmt"
	"gin-bee/async_task/task"
	"gin-bee/config"
	"strings"

	"github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
)

func StartServer() (*machinery.Server, error) {
	cnf := config.Cfg.Machinery
	brokerAddr, brokerPwd := GetRedisAddrAndPwd(cnf.Broker)
	backendAddr, backendPwd := GetRedisAddrAndPwd(cnf.Broker)
	broker := redisbroker.NewGR(cnf, []string{fmt.Sprintf("%s@%s", brokerPwd, brokerAddr)}, 0)
	backend := redisbackend.NewGR(cnf, []string{fmt.Sprintf("%s@%s", backendPwd, backendAddr)}, 0)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)
	// Register tasks
	tasks := map[string]interface{}{
		"TestAdd":        task.Add,
		"UpdateTaskInfo": task.UpdateTaskInfo,
	}
	return server, server.RegisterTasks(tasks)
}

func GetRedisAddrAndPwd(url string) (addr string, passwd string) {
	redisInfo := strings.FieldsFunc(url, func(r rune) bool {
		return r == '@' || r == '/'
	})
	if len(redisInfo) == 2 {
		addr = redisInfo[2]
	} else if len(redisInfo) == 3 {
		addr = redisInfo[2]
		passwd = redisInfo[1]
	}
	return
}
