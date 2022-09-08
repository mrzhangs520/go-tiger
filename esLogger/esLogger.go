package esLogger

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"lib/config"
	"lib/tools"
	"time"
)

// 错误等级
const (
	LeverInfo   = "info"
	LeverWaning = "warning"
	LeverError  = "error"
)

// 日志类型分组
const (
	HttpServer     = "httpServer"
	YgRenderCenter = "ygRenderCenter"
)

var logStashHost = ""
var secret = ""
var source = ""
var mode = ""

type logDataType struct {
	Secret   string `json:"secret"`
	Source   string `json:"source"`
	Mode     string `json:"mode"`
	LogLevel string `json:"log_level"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

// 初始化
func init() {
	logStashHost = config.GetInstance().Section("logStash").Key("host").Value()
	secret = config.GetInstance().Section("logStash").Key("secret").Value()
	source = config.GetInstance().Section("core").Key("serverName").Value()
	mode = tools.Mode
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

	logData := logDataType{
		Secret:   secret,
		Source:   source,
		Mode:     mode,
		LogLevel: logLevel,
		Type:     typeString,
		Message:  messageNew,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}
	logDataJson, err := json.Marshal(logData)
	if nil != err {
		fmt.Printf("日志记录错误：%s\n", err.Error())
	}
	// 调用接口
	resp, err := resty.
		New().
		SetTimeout(10*time.Second).
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(logDataJson).
		Post(logStashHost)
	if nil != err {
		fmt.Printf("日志记录错误：%s\n", err.Error())
	}
	if 200 != resp.StatusCode() || "ok" != string(resp.Body()) {
		fmt.Printf("日志记录错误：%s\n", "接口返回出错！")
	}
}
