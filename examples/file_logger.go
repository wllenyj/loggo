package main

import (
	"github.com/wllenyj/loggo"
)

func main() {
	format := "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}"
	logger := loggo.NewFileLevelLogger(loggo.DEBUG, "debug.log", format)
	loggo.SetDefaultLogger(logger)

	loggo.Debug("Hello world!")

	loggo.Close()
}
