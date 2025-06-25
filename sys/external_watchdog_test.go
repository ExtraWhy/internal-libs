package sys

import (
	"testing"
)

var globCnt = 0

func TestExternalWatchdog(t *testing.T) {
	wd := MakeExternalWatchdog(1)
	defer wd.Stop()
	wd.Start(func() (any, error) {
		t.Log("Why?")
		globCnt++
		return globCnt, nil
	})
}
