// +build !windows

// Copyright 2013, Örjan Persson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loggo

import (
	"fmt"
	"io"
)

type color int

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

var (
	colors = []string{
		UNKNOWN: ColorSeq(ColorWhite),
		FATAL:   ColorSeq(ColorMagenta),
		ERROR:   ColorSeq(ColorRed),
		WARNING: ColorSeq(ColorYellow),
		NOTICE:  ColorSeq(ColorGreen),
		INFO:	 ColorSeq(ColorBlue),
		DEBUG:   ColorSeq(ColorCyan),
	}
	boldcolors = []string{
		UNKNOWN: ColorSeqBold(ColorWhite),
		FATAL:   ColorSeqBold(ColorMagenta),
		ERROR:   ColorSeqBold(ColorRed),
		WARNING: ColorSeqBold(ColorYellow),
		NOTICE:  ColorSeqBold(ColorGreen),
		INFO:	 ColorSeqBold(ColorBlue),
		DEBUG:   ColorSeqBold(ColorCyan),
	}
)

// ConvertColors takes a list of ints representing colors for log levels and
// converts them into strings for ANSI color formatting
func ConvertColors(colors []int, bold bool) []string {
	converted := []string{}
	for _, i := range colors {
		if bold {
			converted = append(converted, ColorSeqBold(color(i)))
		} else {
			converted = append(converted, ColorSeq(color(i)))
		}
	}

	return converted
}

func ColorSeq(color color) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

func ColorSeqBold(color color) string {
	return fmt.Sprintf("\033[%d;1m", int(color))
}

func doFmtVerbLevelColor(layout string, level Flag, output io.Writer) {
	if layout == "bold" {
		output.Write([]byte(boldcolors[level]))
	} else if layout == "reset" {
		output.Write([]byte("\033[0m"))
	} else {
		output.Write([]byte(colors[level]))
	}
}
