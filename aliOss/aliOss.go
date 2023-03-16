package aliOss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mrzhangs520/go-tiger/config"
	"hash"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type myOssType struct {
	bucket *oss.Bucket
}

func New() (*myOssType, error) {
	aliOssConfig := config.GetInstance().Section("aliOss")
	endpoint := aliOssConfig.Key("endpoint").Value()
	accessKeyId := aliOssConfig.Key("accessKeyId").Value()
	accessKeySecret := aliOssConfig.Key("accessKeySecret").Value()
	bucketName := aliOssConfig.Key("bucketName").Value()

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	myOss := new(myOssType)
	myOss.bucket = bucket
	return myOss, nil
}

func (m *myOssType) UploadFile(localFilePath, dir, fileName string, options ...oss.Option) (string, error) {
	// 获取本地文件名
	serverName := config.GetInstance().Section("core").Key("serverName").Value()
	cdnHost := config.GetInstance().Section("aliOss").Key("cdnHost").Value()

	date := time.Now().Format("200601/02")
	unixMicro := time.Now().UnixMicro()

	// 组装上传后的oss地址
	filePath := fmt.Sprintf("%s/%s/%s/%d/%s", serverName, dir, date, unixMicro, fileName)

	// 上传
	err := m.bucket.PutObjectFromFile(filePath, localFilePath, options...)
	if err != nil {
		return "", err
	}

	// 返回新地址
	return fmt.Sprintf("%s/%s", cdnHost, filePath), nil
}

func (m *myOssType) SymlinkFile(originFilePath, dir, fileName string, options ...oss.Option) (string, error) {
	serverName := config.GetInstance().Section("core").Key("serverName").Value()
	cdnHost := config.GetInstance().Section("aliOss").Key("cdnHost").Value()

	originFilePathArr := strings.Split(originFilePath, fmt.Sprintf("%s/", cdnHost))
	originFilePath = originFilePathArr[0]
	if len(originFilePathArr) >= 2 {
		originFilePath = originFilePathArr[1]
	}
	// 组装新地址
	date := time.Now().Format("200601/02")
	unixMicro := time.Now().UnixMicro()
	newFilePath := fmt.Sprintf("%s/%s/%s/%d/%s", serverName, dir, date, unixMicro, fileName)

	// 软连接
	err := m.bucket.PutSymlink(newFilePath, originFilePath, options...)
	if err != nil {
		return "", err
	}

	// 返回新地址
	return fmt.Sprintf("%s/%s", cdnHost, newFilePath), nil
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
func (m *myOssType) GetToken(path string) (policyTokenType, error) {
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
		return policyTokenType{}, err
	}
	deByte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	_, err = io.WriteString(h, deByte)
	if err != nil {
		return policyTokenType{}, err
	}
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken policyTokenType
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expireEndTime
	policyToken.Signature = signedStr
	policyToken.Directory = path
	policyToken.Policy = deByte

	return policyToken, nil
}

// IsFileExist 判断文件是否存在
func (m *myOssType) IsFileExist(path string) bool {
	res, err := m.bucket.IsObjectExist(path)
	if nil != err {
		return false
	}
	return res
}

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

// HandleUrlHost 地址域名转化成内网
func HandleUrlHost(hostUrl string) string {
	aliOssConfig := config.GetInstance().Section("aliOss")
	oldUrl := aliOssConfig.Key("cdnHost").Value()
	newUrl := aliOssConfig.Key("internalHost").Value()
	hostUrl = strings.Replace(hostUrl, oldUrl, newUrl, 1)
	return hostUrl
}

// HandleUrlUnicode 中文处理成unicode
func HandleUrlUnicode(hostUrl string) string {
	re := regexp.MustCompile("[\u4e00-\u9fa5]+")
	chineseChar := re.FindAllString(hostUrl, -1)
	for _, r := range chineseChar {
		hostUrl = strings.Replace(hostUrl, r, url.QueryEscape(r), -1)
	}
	return hostUrl
}
