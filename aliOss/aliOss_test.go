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
