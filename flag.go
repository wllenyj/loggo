package loggo

import (
//"fmt"
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

	LEVEL_ALL

	CUSTOM
)

const (
	LEVEL_MASK = LEVEL_ALL - 1
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
	case level & FATAL == FATAL:
		return "FATAL"
	case level & ERROR == ERROR:
		return "ERROR"
	case level & WARNING == WARNING:
		return "WARNING"
	case level & NOTICE == NOTICE:
		return "NOTICE"
	case level & INFO == INFO:
		return "INFO"
	case level & DEBUG == DEBUG:
		return "DEBUG"
	default:
	}
	return "OTHER"
}
