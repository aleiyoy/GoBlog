package global

import (
	"GoBlog/pkg/logger"
	"GoBlog/pkg/setting"
)



var (
	// 对3个区段配置进行全局变量声明
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS

	// 日志全局
	Logger			*logger.Logger
)