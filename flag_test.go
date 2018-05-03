package loggo

import (
	"testing"
)

func TestFlag(t *testing.T) {
	flags := []struct {
		f Flag
		n string
	}{
		{ FATAL, "FATAL"},
		{ ERROR, "ERROR"},
		{ WARNING, "WARNING"},
		{ NOTICE, "NOTICE"},
		{ INFO, "INFO"},
		{ DEBUG, "DEBUG"},
	}

	for _, f := range flags {
		if f.n != f.f.String() {
			t.Errorf("%b %s != %s", f.f, f.n, f.f)
		}
	}
}

func TestParseFlag(t *testing.T) {
	test_data := []struct {
		f Flag
		n string
	}{
		{ FATAL, "fatal"},
		{ ERROR, "ERROR"},
		{ WARNING, "WARNING"},
		{ WARNING, "Warn"},
		{ NOTICE, "NOTICE"},
		{ INFO, "INFO"},
		{ DEBUG, "Debug"},
		{ DEBUG, "debug"},
		{ DEBUG, "DEBUG"},
	}

	for _, data := range test_data{
		if data.f != ParseLevel(data.n) {
			t.Errorf("%b %s != %s", data.f, data.n, data.f)
		}
	}
}
