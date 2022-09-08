package config

import (
	"fmt"
	"github.com/mrzhangs520/go-tiger/dError"
	"github.com/mrzhangs520/go-tiger/tools"
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
