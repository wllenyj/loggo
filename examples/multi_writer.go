package main

import (
	"github.com/wllenyj/loggo"
	"log"
	"os"
	"os/signal"
	"time"
	"syscall"
)

func main() {
	aaa := 123
	loggo.SetDefaultLogger(loggo.NewFileLevelLogger(loggo.DEBUG, "debug.log", "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}"))
	
	loggo.Debugf("aaa: %d", aaa)
	loggo.Warnf("aaa addr: %p", &aaa)
	loggo.Error("logger %s ", aaa, aaa, " log err.")
	loggo.Fatal("Fatal log err.")
	go func() {
		for {
			time.Sleep(1 * time.Second)
			loggo.Info("logger %s ", aaa, aaa, " log err.")
		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			loggo.Errorf("logger %s ", aaa, aaa, " log err.")
		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			loggo.Warnf("aaa addr: %p", &aaa)
		}
	}()

	//f, _ := os.Create("test.log")
	//stdlog := log.New(f, "", log.LstdFlags)
	stdlog := log.New(os.Stdout, "", log.LstdFlags)
	stdlog.Print("log print", "aaaaa")
	stdlog.Print("log print", "aaaaa")

	//time.Sleep(1 * time.Second)

	flog := loggo.NewFileLevelLogger(loggo.INFO, "info.log", "%{time:01-02 15:04:05.9999} #%{pid}.%{gid} %{shortpkg}/%{shortfile}/%{callpath:3} %{color:bold}%{level:.4s}%{color:reset} %{message}")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			flog.Debugf("aaa: %d", aaa)
			flog.Infof("aaa: %d", aaa)
			flog.Warnf("aaa addr: %p", &aaa)
			flog.Error("logger %s", aaa, aaa, "log err.")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)
	for {
		s := <-c
		switch s {
		case syscall.SIGUSR1:
			err := loggo.Reopen()
			loggo.Info("loggo reopen failed. %v", err)
		case syscall.SIGINT:
			loggo.Info("break")
			goto END_FOR	
		}
	}
	END_FOR:

	loggo.Close()
}
