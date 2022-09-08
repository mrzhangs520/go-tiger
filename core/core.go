package core

import (
	"fmt"
	"os"
	"time"
)

func Start() {
	// 创建temp文件夹
	_ = os.Mkdir(fmt.Sprintf("%s/temp", AppPath), 0777)

	fmt.Printf("当前时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("运行路径：%s\n", AppPath)
	fmt.Printf("运行环境：%s\n", Mode)
}
