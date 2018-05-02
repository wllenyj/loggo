package loggo

import (
	"fmt"
	"strings"
)

type Flag uint32

const (
	ALL Flag = 1 << iota
	FATAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG

	LEVEL_END

	CUSTOM
)

const (
	LEVEL_MASK = LEVEL_END - 1
)

var levelNames = []string{
	"FATAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

//var flagIndex = [...]uint32 {
//	0, 1, 1, 2, 2, 2, 2, 3, 3, 3,
//}

func (flag Flag) String() string {
	level := flag & LEVEL_MASK
	//fmt.Printf("%b %d\n", level, level/2)
	switch {
	case level&FATAL == FATAL:
		return "FATAL"
	case level&ERROR == ERROR:
		return "ERROR"
	case level&WARNING == WARNING:
		return "WARNING"
	case level&NOTICE == NOTICE:
		return "NOTICE"
	case level&INFO == INFO:
		return "INFO"
	case level&DEBUG == DEBUG:
		return "DEBUG"
	default:
	}
	return "OTHER"
}

func ParseLevel(lvl string) Flag {
	switch strings.ToLower(lvl) {
	case "fatal":
		return FATAL
	case "error":
		return ERROR
	case "warn", "warning":
		return WARNING
	case "notice", "noti":
		return NOTICE
	case "info":
		return INFO
	case "debug":
		return DEBUG
	default:
		panic(fmt.Errorf("not a valid Level: %q", lvl))
	}
	return ALL
}
