package logger

import (
	"sync"
	"yunyandz.com/tiktok/user-part/settings"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	suger  *zap.SugaredLogger
	once   sync.Once
)

// 创建新的logger
func New() *zap.Logger {
	once.Do(func() {
		if !settings.ServiceConf.Debug {
			logger, _ = zap.NewProduction()
		} else {
			logger, _ = zap.NewDevelopment()
			logger.Debug("running in debug mode...")
		}
		suger = logger.Sugar()
		zap.ReplaceGlobals(logger)
	})
	return logger
}

// 返回一个suger模式下的logger，可以引入后直接使用Suger().xxx()
func Suger() *zap.SugaredLogger {
	return suger
}

func Logger() *zap.Logger {
	return logger
}
