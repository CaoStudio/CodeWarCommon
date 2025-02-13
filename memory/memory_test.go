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
		t.Log(humanize.Bytes(Size(pointStr)), " ", len(str))
	}
}

// TestMap test map
func TestMap(t *testing.T) {
	m := make(map[int]int)
	for i := 0; i < 1000; i++ {
		m[i] = i
		t.Log(humanize.Bytes(Size(m)))
	}
}
