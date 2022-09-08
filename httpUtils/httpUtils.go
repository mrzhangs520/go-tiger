package httpUtils

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/mrzhangs520/go-tiger/esLogger"
	"net/http"
)

func Success(c *gin.Context, data ...interface{}) {
	var res interface{}
	res = nil
	if 0 != len(data) {
		res = data[0]
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":   http.StatusOK,
			"msg":    "ok!",
			"result": res,
		},
	)
}

func Error(c *gin.Context, err string) {
	c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"code":   http.StatusInternalServerError,
			"msg":    err,
			"result": nil,
		},
	)
	input, _ := c.Get("param")

	// 记录错误日志
	msg := map[string]interface{}{
		"requestUrl":  c.Request.URL.RequestURI(),
		"requestBody": input,
		"errorMsg":    err,
	}
	esLogger.Write(esLogger.LeverError, esLogger.HttpServer, msg)
}

// GetBodyParam map转结构体
func GetBodyParam(c *gin.Context, output interface{}) {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}
	// 忽略解析map到结构体的报错，因为控制器会做验证
	decoder, _ := mapstructure.NewDecoder(config)
	input, _ := c.Get("param")
	_ = decoder.Decode(input)
}
