package global

import (
	"github.com/fishblog/pkg/logger"
	"github.com/fishblog/pkg/setting"
)

var (
	ServerSetting  *setting.ServerSettings
	AppSetting		*setting.AppSettings
	DatabaseSetting *setting.DatabaseSettings
	Logger 		    *logger.Logger
)


