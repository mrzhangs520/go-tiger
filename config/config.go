package config

import (
	"fmt"
	"go-tiger/dError"
	"go-tiger/tools"
	"gopkg.in/ini.v1"
)

var instance *ini.File

func init() {
	configPath := fmt.Sprintf("%s/conf/%s.ini", tools.AppPath, tools.Mode)
	var err error
	instance, err = ini.Load(configPath)
	if err != nil {
		panic(dError.NewError("config.GetInstance.ini.Load", err))
	}
}

func GetInstance() *ini.File {
	return instance
}
