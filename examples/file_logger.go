package main

import (
	"github.com/wllenyj/loggo"
)

func main() {
	logger := loggo.NewFileLevelLogger(loggo.DEBUG, "debug.log", loggo.DefFormat)
	loggo.SetDefaultLogger(logger)

	loggo.Debug("Hello world!")

	loggo.Close()
}
