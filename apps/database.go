package apps

import (
	"gin-bee/config"
)

var Db, _ = config.DB(*config.Cfg.Database)
