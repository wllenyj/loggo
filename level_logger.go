package loggo

import (
	"time"
	"os"
	"github.com/petermattis/goid"
)

type levelLogger struct {
	level Flag
	
	formatter Formatter 

	writer stringWriter
}

func (lg *levelLogger) Debug(args ...interface{}) {
	lg.log(DEBUG, 1, nil, args...)
}
func (lg *levelLogger) Debugf(format string, args ...interface{}) {
	lg.log(DEBUG, 1, &format, args...)
}
func (lg *levelLogger) Info(args ...interface{}) {
	lg.log(INFO, 1, nil, args...)
}
func (lg *levelLogger) Infof(format string, args ...interface{}) {
	lg.log(INFO, 1, &format, args...)
}
func (lg *levelLogger) Warn(args ...interface{}) {
	lg.log(WARNING, 1, nil, args...)
}
func (lg *levelLogger) Warnf(format string, args ...interface{}) {
	lg.log(WARNING, 1, &format, args...)
}
func (lg *levelLogger) Error(args ...interface{}) {
	lg.log(ERROR, 1, nil, args...)
}
func (lg *levelLogger) Errorf(format string, args ...interface{}) {
	lg.log(ERROR, 1, &format, args...)
}
func (lg *levelLogger) Fatal(args ...interface{}) {
	lg.log(FATAL, 1, nil, args...)
}
func (lg *levelLogger) Fatalf(format string, args ...interface{}) {
	lg.log(FATAL, 1, &format, args...)
}

func (lg *levelLogger) Output(calldepth int, s string) error {
	lg.log(ALL, 1, &s, nil)
	return nil
}

func (lg *levelLogger) log(level Flag, calldepth int, s *string, args ...interface{} ) {
	if lg.level < level {
		return
	}

	r := GetRecord()
	r.level = level
	r.gid = goid.Get()
	r.time = time.Now()
	r.fmt = s
	r.args = args

	str := FormatterProxy(lg.formatter, calldepth+1, r)
	PutRecord(r)
	lg.writer.WriteString(str)
}

//func (lg *levelLogger) SetFileWriter(filename string) {
//	lg.writer = GetFileWriter(filename)
//}
//func (lg *levelLogger) SetFormatter(format string) {
//	lg.formatter = MustStringFormatter(fmt)
//}
//func (lg *levelLogger) SetLevel(level Flag) {
//	lg.level = level
//}

func NewLevelLogger(level Flag, fmt string) LevelLogger {
	return &levelLogger {
		level: level,
		writer: os.Stdout,
		formatter: MustStringFormatter(fmt),
	}
}
func NewFileLevelLogger(level Flag, filename string, fmt string) LevelLogger {
	return &levelLogger {
		level: level,
		writer: GetFileWriter(filename),
		formatter: MustStringFormatter(fmt),
	}
}
