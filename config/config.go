package config

import (
	"fmt"
	"github.com/mrzhangs520/go-tiger/core"
	"github.com/mrzhangs520/go-tiger/dError"
	"gopkg.in/ini.v1"
)

var instance *ini.File

func init() {
	configPath := fmt.Sprintf("%s/conf/%s.ini", core.AppPath, core.Mode)
	var err error
	instance, err = ini.Load(configPath)
	if err != nil {
		panic(dError.NewError("读取配置文件出错", err))
	}
}

func GetInstance() *ini.File {
	return instance
}
