package aliOss

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mrzhangs520/go-tiger/core"
	"testing"
)

func TestGetToken(t *testing.T) {
	core.Start()
	token := New().GetToken("/go-api-shop")
	spew.Dump(token)
}

func TestIsFileExist(t *testing.T) {
	core.Start()
	res := New().isFileExist("go-api-shop/system/cut-img/go-api-shop/upload/media/origin-mov/of1r9-e1df-fq1df/测试.jpg")
	spew.Dump(res)
}
