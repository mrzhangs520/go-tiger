package aliOss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mrzhangs520/go-tiger/config"
	"github.com/mrzhangs520/go-tiger/dError"
	"hash"
	"io"
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

type policyTokenType struct {
	AccessKeyId string `json:"access_id"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

// GetToken 获取token
func (m *myOssType) GetToken(path string) policyTokenType {
	var err error

	aliOssConfig := config.GetInstance().Section("aliOss")
	accessKeyId := aliOssConfig.Key("accessKeyId").Value()
	accessKeySecret := aliOssConfig.Key("accessKeySecret").Value()
	host := aliOssConfig.Key("cdnHost").Value()

	type ConfigStructType struct {
		Expiration string     `json:"expiration"`
		Conditions [][]string `json:"conditions"`
	}

	now := time.Now().Unix()
	// 过期时间600s
	expireEndTime := now + 3600
	var tokenExpire = getGmtIso8601(expireEndTime)

	//create post policy json
	configInfo := ConfigStructType{}
	configInfo.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, path)
	configInfo.Conditions = append(configInfo.Conditions, condition)

	// 计算签名
	result, err := json.Marshal(configInfo)
	if nil != err {
		panic(dError.NewError("获取上传签名错误", "aliOss.GetToken.json.Marshal(configInfo)", err))
	}
	deByte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	_, err = io.WriteString(h, deByte)
	if err != nil {
		panic(dError.NewError("获取上传签名错误", "aliOss.GetToken.io.WriteString", err))
	}
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken policyTokenType
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expireEndTime
	policyToken.Signature = signedStr
	policyToken.Directory = path
	policyToken.Policy = deByte

	return policyToken
}

// IsFileExist 判断文件是否存在
func (m *myOssType) IsFileExist(path string) bool {
	res, err := m.bucket.IsObjectExist(path)
	if nil != err {
		panic(dError.NewError("上传系统错误", "aliOss.GetToken.bucket.SignURL", err))
	}
	return res
}

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}
