package dLogger

import (
	"github.com/mrzhangs520/go-tiger/core"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	core.Start()
	Write(LeverWaning, "dfdsf", "fdf")
	time.Sleep(time.Second)
}
