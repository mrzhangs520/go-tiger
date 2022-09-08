package aliOss

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestHandleUrlHost(t *testing.T) {
	spew.Dump(HandleUrlHost("https://produce01.xiyin.love/newbox/static/icon/飞机盒视图/上面.png"))
}
