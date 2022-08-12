package zaplog

import (
	"go.uber.org/zap"
)

var Logger = GetLogger()

func GetLogger() *zap.Logger {
	lg, _ := zap.NewDevelopment()
	return lg
}
