package aliOss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mrzhangs520/go-tiger/config"
	"github.com/mrzhangs520/go-tiger/dError"
	"path/filepath"
	"time"
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
		panic(dError.NewError("上传系统错误", "aliOss.UploadFile.oss.New", err))
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic(dError.NewError("上传系统错误", "aliOss.UploadFile.client.Bucket", err))
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
		panic(dError.NewError("上传系统错误", "aliOss.UploadFile.bucket.PutObjectFromFile", err))
	}

	// 返回新地址
	return fmt.Sprintf("%s/%s", cdnHost, filePath)
}

// GetToken 获取token
func (m *myOssType) GetToken(path string) string {
	token, err := m.bucket.SignURL(path, oss.HTTPPut, 600)
	if nil != err {
		panic(dError.NewError("上传系统错误", "aliOss.GetToken.bucket.SignURL", err))
	}
	return token
}

// IsFileExist 判断文件是否存在
func (m *myOssType) IsFileExist(path string) bool {
	res, err := m.bucket.IsObjectExist(path)
	if nil != err {
		panic(dError.NewError("上传系统错误", "aliOss.GetToken.bucket.SignURL", err))
	}
	return res
}
