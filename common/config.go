package common

import (
	"github.com/zhao890420/goconfig"
	"go.uber.org/zap"
	"time"
)

var config *goconfig.ConfigFile
var DefLogger *zap.Logger
var AccessLogger *zap.Logger
var TimeLocation *time.Location

type Config struct {
}

func GetConfig() *goconfig.ConfigFile {
	return config
}
func InitConfigFile(configFile string) {
	var err error
	config, err = goconfig.LoadConfigFile(configFile)
	if err != nil {
		panic(err)
	}
	// 设置时区
	if location, err := time.LoadLocation("Asia/Chongqing"); err != nil {
		panic(err)
	} else {
		TimeLocation = location
	}
	DefLogger.Info("======finish configFile init ")
}

func (c Config) Destroy() {
}
