package esLogger

// 错误等级
const (
	LeverInfo   = "info"
	LeverWaning = "warning"
	LeverError  = "error"
)

type logDataType struct {
	Secret   string `json:"secret"`
	Source   string `json:"source"`
	Mode     string `json:"mode"`
	LogLevel string `json:"log_level"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

// Write 写入日志到logStash
func Write(logLevel, typeString string, message interface{}) {
	// 开启一个携程异步写入
	go func() {
		toWrite(logLevel, typeString, message)
	}()
}

func toWrite(logLevel, typeString string, message interface{}) {

}
