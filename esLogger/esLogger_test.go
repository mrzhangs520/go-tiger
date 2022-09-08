package esLogger

import (
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	msg := map[string]string{
		"测试":  "狮子测试一下日志写入",
		"测试1": "狮子测试一下日志写入1",
	}
	Write(LeverInfo, HttpServer, msg)

	time.Sleep(time.Second * 3)
}
