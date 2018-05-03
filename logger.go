package loggo

import ()

type Logger interface {
	log(Flag, int, *string, ...interface{})
}

type LevelLogger interface {
	Logger

	StdOutput

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type StdOutput interface {
	Output(calldepth int, s string) error
}
