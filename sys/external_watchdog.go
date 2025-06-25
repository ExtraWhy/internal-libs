package sys

import (
	"fmt"
	"time"
)

type ExternalWatchdog struct {
	tick *time.Ticker
}

func MakeExternalWatchdog(interval int) *ExternalWatchdog {
	ew := &ExternalWatchdog{
		tick: time.NewTicker(time.Duration(interval) * time.Second),
	}
	return ew
}

func (ew *ExternalWatchdog) Start(cb func() (any, error)) {
	if cb != nil {
		go func() {
			for t := range ew.tick.C {
				fmt.Println("tick at ", t)
				fmt.Println(cb())
			}
		}()
	}
}

func (ew *ExternalWatchdog) Stop() {
	ew.tick.Stop()
}

func (ew *ExternalWatchdog) Restart() {

}
