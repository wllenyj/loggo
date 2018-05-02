package loggo

import ()

var (
	default_logger Logger

	DefFormat = "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}"
)

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

func SetDefaultLogger(logger Logger) {
	default_logger = logger
}

func init() {
	default_logger = NewLevelLogger(INFO, DefFormat)
	//default_logger = NewFileLevelLogger(INFO,"debug.log", "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}")
}
