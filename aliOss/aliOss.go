package aliOss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mrzhangs520/go-tiger/config"
	"github.com/mrzhangs520/go-tiger/dError"
	"github.com/mrzhangs520/go-tiger/tools"
	"path/filepath"
	"strings"
	"time"
)

// 上传的文件路径
const (
	SurfaceImage = "surface"
	TextImage    = "textImage"
	AppendGraph  = "appendGraph"
)

type myOssType struct {
	bucket *oss.Bucket
}

func New() *myOssType {
	aliOssConfig := config.GetInstance().Section("aliOss")
	endpoint := aliOssConfig.Key("endpoint").Value()
	accessKeyId := aliOssConfig.Key("accessKeyId").Value()
	accessKeySecret := aliOssConfig.Key("accessKeySecret").Value()
	bucketName := aliOssConfig.Key("bucketName").Value()

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		panic(dError.NewError("aliOss.UploadFile.oss.New", err))
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic(dError.NewError("aliOss.UploadFile.client.Bucket", err))
	}
	myOss := new(myOssType)
	myOss.bucket = bucket
	return myOss
}

func (m *myOssType) UploadFile(localFilePath, dir string) string {
	// 获取本地文件名
	_, fileName := filepath.Split(localFilePath)
	serverName := config.GetInstance().Section("core").Key("serverName").Value()
	cdnHost := config.GetInstance().Section("aliOss").Key("cdnHost").Value()

	// 拼接上项目地址
	date := time.Now().Format("200601/02")
	filePath := fmt.Sprintf("%s/%s/%s/%s", serverName, dir, date, fileName)

	// 上传
	err := m.bucket.PutObjectFromFile(filePath, localFilePath)
	if err != nil {
		panic(dError.NewError("aliOss.UploadFile.bucket.PutObjectFromFile", err))
	}

	// 返回新地址
	return fmt.Sprintf("%s/%s", cdnHost, filePath)
}

// HandleUrlHost 地址域名转化成内网
func HandleUrlHost(url string) string {
	if "local" == tools.Mode {
		return url
	}

	url = strings.Replace(url, "produce01.xiyin.love", "baiyin-produce.oss-cn-hangzhou-internal.aliyuncs.com", 1)
	url = strings.Replace(url, "baiyin-produce.oss-cn-hangzhou.aliyuncs.com", "baiyin-produce.oss-cn-hangzhou-internal.aliyuncs.com", 1)
	url = strings.Replace(url, "baiyin-dev.oss-cn-hangzhou.aliyuncs.com", "baiyin-dev.oss-cn-hangzhou-internal.aliyuncs.com", 1)
	url = strings.Replace(url, "dev01.xiyin.love", "baiyin-dev.oss-cn-hangzhou-internal.aliyuncs.com", 1)

	return url
}
