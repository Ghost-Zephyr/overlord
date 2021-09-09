package log

import (
	"fmt"
	"log"
)

var CurrentLogLevel LogLevel

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func (lvl LogLevel) String() string {
	return [...]string{"TRACE", "DEBUG",
		"INFO", "WARN", "ERROR", "FATAL"}[lvl]
}

func (lvl LogLevel) Integer() int {
	return [...]int{0, 1, 2, 3, 4, 5}[lvl]
}

func PrintLog(lvl LogLevel, str string, a ...interface{}) {
	if CurrentLogLevel.Integer() > lvl.Integer() {
		return
	}
	str = fmt.Sprintf(
		"%s %s",
		lvl.String(),
		str)
	if lvl == FATAL {
		log.Fatalf(str, a...)
	} else {
		log.Printf(str, a...)
	}
}
