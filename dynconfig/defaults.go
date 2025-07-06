package dynconfig

import (
	"fmt"
	"time"
)

//------------------------------------------------------------------------

func Insert(a []string, c string, i int) []string {
	return append(a[:i], append([]string{c}, a[i:]...)...)
}

//------------------------------------------------------------------------

type defaultMetrics struct{}

func (m *defaultMetrics) Timing(metricName string, duration time.Duration, tags map[string]string) {
	if tags == nil {
		tags = make(map[string]string)
	}

	currentTime := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] Timing metric: %s, Duration: %s, Tags: %v\n", currentTime, metricName, duration, tags)
}

func (m *defaultMetrics) Increment(metricName string, value int, tags map[string]string) {
	if tags == nil {
		tags = make(map[string]string)
	}
	currentTime := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] Increment metric: %s, Value: %d, Tags: %v\n", currentTime, metricName, value, tags)
}

//------------------------------------------------------------------------

type defaultLogger struct {
}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	l.logMessage("INFO", l.argsToStrings(msg, args))
}

func (l *defaultLogger) Warn(msg string, args ...interface{}) {
	l.logMessage("WARN", l.argsToStrings(msg, args))
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	l.logMessage("ERROR", l.argsToStrings(msg, args))
}

func (l *defaultLogger) logMessage(level string, args []string) {
	currentTime := time.Now().Format(time.RFC3339)

	if len(args)%2 != 0 {
		panic("Arguments must be in key-value pairs")
	}

	formattedArgs := ""
	for i := 0; i < len(args); i += 2 {
		formattedArgs += fmt.Sprintf("%s:%s, ", args[i], args[i+1])
	}

	if len(formattedArgs) > 0 {
		formattedArgs = formattedArgs[:len(formattedArgs)-2] // remove the last comma and space
	}
	fmt.Printf("[%s] %s: %s\n", currentTime, level, formattedArgs)
}

func (l *defaultLogger) argsToStrings(msg string, args []interface{}) []string {
	if len(args) == 0 {
		return []string{}
	}

	var strArgs []string
	var pos int

	if len(args)%2 == 0 {
		strArgs = make([]string, len(args)+2)
		strArgs[0] = "message"
		strArgs[1] = msg
		pos = 2
	} else {
		strArgs = make([]string, len(args)+1)
		strArgs[0] = msg
		pos = 1
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			strArgs[pos] = v
		default:
			strArgs[pos] = fmt.Sprintf("%v", v)
		}
		pos++
	}
	return strArgs
}

//------------------------------------------------------------------------

type defaultScheduler struct{}

func (ds *defaultScheduler) RunAfter(interval time.Duration, f func()) {
	time.AfterFunc(interval, f)
}

//------------------------------------------------------------------------
