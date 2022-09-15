package dLogger

import (
	"encoding/json"
	"fmt"
	"github.com/mrzhangs520/go-tiger/config"
	"github.com/mrzhangs520/go-tiger/core"
	"time"
)

// 错误等级
const (
	LeverInfo   = "info"
	LeverWaning = "warning"
	LeverError  = "error"
)

var source string
var mode string

// 初始化
func init() {
	source = config.GetInstance().Section("core").Key("serverName").Value()
	mode = core.Mode
}

// Write 写入日志到logStash
func Write(logLevel, typeString string, message interface{}) {
	// 开启一个携程异步写入
	go func() {
		toWrite(logLevel, typeString, message)
	}()
}

func toWrite(logLevel, typeString string, message interface{}) {
	// 将message转为字符串
	messageNew, ok := message.(string)
	if !ok {
		messageJson, err := json.Marshal(message)
		if nil != err {
			fmt.Printf("日志记录错误：%s\n", err.Error())
		}
		messageNew = string(messageJson)
	}

	logData := &LogModelType{
		Source:     source,
		Mode:       mode,
		LogLevel:   logLevel,
		Type:       typeString,
		Message:    messageNew,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	_, err := logData.Create()

	if nil != err {
		fmt.Printf("写入日志失败： %s", err.Error())
	}
}
