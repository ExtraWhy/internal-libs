package sys

import (
	"testing"
)

var globCnt = 0

func TestExternalWatchdog(t *testing.T) {
	wd := MakeExternalWatchdog(1)
	defer wd.Stop()
	wd.Start(func() (any, error) {
		globCnt++
		return globCnt, nil
	})
}

func TestExternalWatchdogs(t *testing.T) {
	wd1 := MakeExternalWatchdog(1)
	wd2 := MakeExternalWatchdog(2)
	wd3 := MakeExternalWatchdog(3)
	defer wd1.Stop()
	defer wd2.Stop()
	defer wd3.Stop()
	wd1.Start(func() (any, error) {
		globCnt++
		return globCnt, nil
	})
	wd2.Start(func() (any, error) {
		return 42, nil
	})
	wd3.Start(func() (any, error) {

		return "NOP", nil
	})

	select {}
}
