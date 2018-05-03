// Copyright 2013, Örjan Persson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loggo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// TODO see Formatter interface in fmt/print.go
// TODO try text/template, maybe it have enough performance
// TODO other template systems?
// TODO make it possible to specify formats per backend?
type fmtVerb int

const (
	fmtVerbTime fmtVerb = iota
	fmtVerbLevel
	fmtVerbPid
	fmtVerbGid
	fmtVerbProgram
	fmtVerbMessage
	fmtVerbLongfile
	fmtVerbShortfile
	fmtVerbLongpkg
	fmtVerbShortpkg
	fmtVerbLongfunc
	fmtVerbShortfunc
	fmtVerbCallpath
	fmtVerbLevelColor

	// Keep last, there are no match for these below.
	fmtVerbUnknown
	fmtVerbStatic
)

var fmtVerbs = []string{
	"time",
	"level",
	"pid",
	"gid",
	"program",
	"message",
	"longfile",
	"shortfile",
	"longpkg",
	"shortpkg",
	"longfunc",
	"shortfunc",
	"callpath",
	"color",
}

const rfc3339Milli = "2006-01-02T15:04:05.999999"

var defaultVerbsLayout = []string{
	rfc3339Milli,
	"s",
	"d",
	"d",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"0",
	"",
}

var (
	pid     = os.Getpid()
	program = filepath.Base(os.Args[0])
)

func getFmtVerbByName(name string) fmtVerb {
	for i, verb := range fmtVerbs {
		if name == verb {
			return fmtVerb(i)
		}
	}
	return fmtVerbUnknown
}

type stringWriter interface {
	io.Writer
	WriteString(string) (int, error)
	//WriteByte(c byte) error
}

// Formatter is the required interface for a custom log record formatter.
type Formatter interface {
	Format(calldepth int, r *Record, output stringWriter) error
}

func FormatterProxy(formatter Formatter, calldepth int, r *Record) []byte {
	output := buffer_pool.Get().(*bytes.Buffer)
	output.Reset()
	formatter.Format(calldepth+1, r, output)
	ret := output.Bytes()
	buffer_pool.Put(output)
	return ret
	//return output.String()
}

var (
	// DefaultFormatter is the default formatter used and is only the message.
	DefaultFormatter = MustStringFormatter("%{message}")

	// GlogFormatter mimics the glog format
	GlogFormatter = MustStringFormatter("%{level:.1s}%{time:01-02T15:04:05.999999} %{pid}.%{gid} %{shortfile}] %{message}")

	//builder_pool = sync.Pool{
	//	New: func() interface{} {
	//		return &strings.Builder{}
	//	},
	//}
	buffer_pool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

var formatRe = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`)

type part struct {
	verb   fmtVerb
	layout string
}

// stringFormatter contains a list of parts which explains how to build the
// formatted string passed on to the logging backend.
type stringFormatter struct {
	parts []part
}

// NewStringFormatter returns a new Formatter which outputs the log record as a
// string based on the 'verbs' specified in the format string.
//
// The verbs:
//
// General:
//     %{pid}       Process id (int)
//     %{gid}       Goroutine id (int)
//     %{time}      Time when log occurred (time.Time)
//     %{level}     Log level (Level)
//     %{module}    Module (string)
//     %{program}   Basename of os.Args[0] (string)
//     %{message}   Message (string)
//     %{longfile}  Full file name and line number: /a/b/c/d.go:23
//     %{shortfile} Final file name element and line number: d.go:23
//     %{callpath}  Callpath like main.a.b.c...c  "..." meaning recursive call ~. meaning truncated path
//     %{color}     ANSI color based on log level
//
// For normal types, the output can be customized by using the 'verbs' defined
// in the fmt package, eg. '%{id:04d}' to make the id output be '%04d' as the
// format string.
//
// For time.Time, use the same layout as time.Format to change the time format
// when output, eg "2006-01-02T15:04:05.999Z-07:00".
//
// For the 'color' verb, the output can be adjusted to either use bold colors,
// i.e., '%{color:bold}' or to reset the ANSI attributes, i.e.,
// '%{color:reset}' Note that if you use the color verb explicitly, be sure to
// reset it or else the color state will persist past your log message.  e.g.,
// "%{color:bold}%{time:15:04:05} %{level:-8s}%{color:reset} %{message}" will
// just colorize the time and level, leaving the message uncolored.
//
// For the 'callpath' verb, the output can be adjusted to limit the printing
// the stack depth. i.e. '%{callpath:3}' will print '~.a.b.c'
//
// Colors on Windows is unfortunately not supported right now and is currently
// a no-op.
//
// There's also a couple of experimental 'verbs'. These are exposed to get
// feedback and needs a bit of tinkering. Hence, they might change in the
// future.
//
// Experimental:
//     %{longpkg}   Full package path, eg. github.com/go-logging
//     %{shortpkg}  Base package path, eg. go-logging
//     %{longfunc}  Full function name, eg. littleEndian.PutUint32
//     %{shortfunc} Base function name, eg. PutUint32
//     %{callpath}  Call function path, eg. main.a.b.c
func NewStringFormatter(format string) (Formatter, error) {
	var fmter = &stringFormatter{}

	// Find the boundaries of all %{vars}
	matches := formatRe.FindAllStringSubmatchIndex(format, -1)
	if matches == nil {
		return nil, errors.New("invalid log format: " + format)
	}

	// Collect all variables and static text for the format
	prev := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		if start > prev {
			fmter.add(fmtVerbStatic, format[prev:start])
		}

		name := format[m[2]:m[3]]
		verb := getFmtVerbByName(name)
		if verb == fmtVerbUnknown {
			return nil, errors.New("unknown variable: " + name)
		}

		// Handle layout customizations or use the default. If this is not for the
		// time, color formatting or callpath, we need to prefix with %.
		layout := defaultVerbsLayout[verb]
		if m[4] != -1 {
			layout = format[m[4]:m[5]]
		}
		if verb != fmtVerbTime && verb != fmtVerbLevelColor && verb != fmtVerbCallpath {
			layout = "%" + layout
		}

		fmter.add(verb, layout)
		prev = end
	}
	end := format[prev:]
	if end != "" {
		fmter.add(fmtVerbStatic, end)
	}

	// Make a test run to make sure we can format it correctly.
	//t, err := time.Parse(time.RFC3339, "2010-02-04T21:00:57-08:00")
	//if err != nil {
	//	panic(err)
	//}
	//testFmt := "hello %s"
	//r := &Record{
	//	//ID:     12345,
	//	Time:   t,
	//	Module: "logger",
	//	Args:   []interface{}{"go"},
	//	fmt:    &testFmt,
	//}
	//if err := fmter.Format(0, r, &bytes.Buffer{}); err != nil {
	//	return nil, err
	//}

	return fmter, nil
}

// MustStringFormatter is equivalent to NewStringFormatter with a call to panic
// on error.
func MustStringFormatter(format string) Formatter {
	f, err := NewStringFormatter(format)
	if err != nil {
		panic("Failed to initialized string formatter: " + err.Error())
	}
	return f
}

func (f *stringFormatter) add(verb fmtVerb, layout string) {
	f.parts = append(f.parts, part{verb, layout})
}

func (f *stringFormatter) Format(calldepth int, r *Record, output stringWriter) error {
	for _, part := range f.parts {
		if part.verb == fmtVerbStatic {
			output.WriteString(part.layout)
		} else if part.verb == fmtVerbTime {
			output.WriteString(r.time.Format(part.layout))
		} else if part.verb == fmtVerbLevelColor {
			doFmtVerbLevelColor(part.layout, r.level, output)
		} else if part.verb == fmtVerbCallpath {
			depth, err := strconv.Atoi(part.layout)
			if err != nil {
				depth = 0
			}
			output.WriteString(formatCallpath(calldepth+1, depth))
		} else if part.verb == fmtVerbMessage {
			doFormatMessage(r, output)
		} else {
			var v interface{}
			switch part.verb {
			case fmtVerbLevel:
				v = r.level
				break
			case fmtVerbPid:
				v = pid
				break
			case fmtVerbGid:
				v = r.gid
				break
			case fmtVerbProgram:
				v = program
				break
			case fmtVerbLongfile, fmtVerbShortfile:
				_, file, line, ok := runtime.Caller(calldepth + 1)
				if !ok {
					file = "???"
					line = 0
				} else if part.verb == fmtVerbShortfile {
					file = filepath.Base(file)
				}
				v = fmt.Sprintf("%s:%d", file, line)
			case fmtVerbLongfunc, fmtVerbShortfunc,
				fmtVerbLongpkg, fmtVerbShortpkg:
				// TODO cache pc
				v = "???"
				if pc, _, _, ok := runtime.Caller(calldepth + 1); ok {
					if f := runtime.FuncForPC(pc); f != nil {
						v = formatFuncName(part.verb, f.Name())
					}
				}
			default:
				panic("unhandled format part")
			}
			fmt.Fprintf(output, part.layout, v)
		}
	}
	output.Write([]byte{'\n'})
	return nil
}

func doFormatMessage(r *Record, output stringWriter) {
	if r.fmt != nil {
		fmt.Fprintf(output, *r.fmt, r.args...)
	} else {
		fmt.Fprint(output, r.args...)
	}	
}

// formatFuncName tries to extract certain part of the runtime formatted
// function name to some pre-defined variation.
//
// This function is known to not work properly if the package path or name
// contains a dot.
func formatFuncName(v fmtVerb, f string) string {
	i := strings.LastIndex(f, "/")
	j := strings.Index(f[i+1:], ".")
	if j < 1 {
		return "???"
	}
	pkg, fun := f[:i+j+1], f[i+j+2:]
	switch v {
	case fmtVerbLongpkg:
		return pkg
	case fmtVerbShortpkg:
		return path.Base(pkg)
	case fmtVerbLongfunc:
		return fun
	case fmtVerbShortfunc:
		i = strings.LastIndex(fun, ".")
		return fun[i+1:]
	}
	panic("unexpected func formatter")
}

func formatCallpath(calldepth int, depth int) string {
	v := ""
	callers := make([]uintptr, 64)
	n := runtime.Callers(calldepth+2, callers)
	oldPc := callers[n-1]

	start := n - 3
	if depth > 0 && start >= depth {
		start = depth - 1
		v += "~."
	}
	recursiveCall := false
	for i := start; i >= 0; i-- {
		pc := callers[i]
		if oldPc == pc {
			recursiveCall = true
			continue
		}
		oldPc = pc
		if recursiveCall {
			recursiveCall = false
			v += ".."
		}
		if i < start {
			v += "."
		}
		if f := runtime.FuncForPC(pc); f != nil {
			v += formatFuncName(fmtVerbShortfunc, f.Name())
		}
	}
	return v
}
