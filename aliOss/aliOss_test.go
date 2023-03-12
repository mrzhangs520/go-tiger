package aliOss

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mrzhangs520/go-tiger/core"
	"testing"
)

func TestHandleUrlHost(t *testing.T) {
	core.Start()
	url := "https://oss01.tiger12345.cc/go-api-shop/upload/2023-03-12/中文重重.JPG"
	url = HandleUrlUnicode(url)
	spew.Dump(url)
	url = HandleUrlHost(url)
	spew.Dump(url)
}
