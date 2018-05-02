package loggo

import ()

var (
	default_logger Logger
)

type loggo struct {
	//backend Backend
}

func (log *loggo) Log() {
	//log.backend.Log()
}

//func SetBackend(be ...Backend) {
//	if len(be) == 1 {
//
//	} else {
//
//	}
//	return
//}
//func Close() {
//	default_logger
//}

func Debug(args ...interface{}) {
	default_logger.log(DEBUG, 1, nil, args...)
}
func Debugf(format string, args ...interface{}) {
	default_logger.log(DEBUG, 1, &format, args...)
}
func Info(args ...interface{}) {
	default_logger.log(INFO, 1, nil, args...)
}
func Infof(format string, args ...interface{}) {
	default_logger.log(INFO, 1, &format, args...)
}
func Warn(args ...interface{}) {
	default_logger.log(WARNING, 1, nil, args...)
}
func Warnf(format string, args ...interface{}) {
	default_logger.log(WARNING, 1, &format, args...)
}
func Error(args ...interface{}) {
	default_logger.log(ERROR, 1, nil, args...)
}
func Errorf(format string, args ...interface{}) {
	default_logger.log(ERROR, 1, &format, args...)
}
func Fatal(args ...interface{}) {
	default_logger.log(FATAL, 1, nil, args...)
}
func Fatalf(format string, args ...interface{}) {
	default_logger.log(FATAL, 1, &format, args...)
}

//func SetFileWriter(filename string) {
//	default_logger.SetFileWriter(filename)
//}
//func SetFormatter(format string) {
//	default_logger.SetFormatter(format)
//}
//func SetLevel(level Flag) {
//	default_logger.SetLevel(level)
//}
func SetDefaultLogger(logger Logger) {
	default_logger = logger
}

func init() {
	default_logger = NewLevelLogger(INFO, "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}")
	//default_logger = NewFileLevelLogger(INFO,"debug.log", "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}")
}
