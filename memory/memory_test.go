package memory

import (
	"github.com/dustin/go-humanize"
	"testing"
)

func TestCount(t *testing.T) {
	var pointStr *string
	str := "hello"
	pointStr = &str
	for i := 0; i < 1000; i++ {
		*pointStr = *pointStr + "pointStr"
		t.Log(humanize.Bytes(uint64(Size(pointStr))), " ", len(str))
	}
}
