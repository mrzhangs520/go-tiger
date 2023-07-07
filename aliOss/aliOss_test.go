package aliOss

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mrzhangs520/go-tiger/core"
	"testing"
)

func TestHandleUrlHost(t *testing.T) {
	core.Start()
	url := "https://oss01.tiger12345.cc/go-api-shop/upload/2023-03-12/中文重重.JPG"
	spew.Dump(url)
	url = HandleUrlHost(url)
	spew.Dump(url)
}

func Test_SymlinkFile(t *testing.T) {
	core.Start()
	url := "https://oss01.tiger12345.cc/go-api-supply-chain/upload/2023-03-15/1678887430469-aJiBk7-qZysch/template.png"
	oss, _ := New()

	spew.Dump(url)
	spew.Dump(oss.SymlinkFile(url, "symlink", "tedddst.png"))
}
