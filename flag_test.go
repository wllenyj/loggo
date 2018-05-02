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
		{ INFO|DEBUG, "INFO"},
	}

	for _, f := range flags {
		if f.n != f.f.String() {
			t.Errorf("%b %s != %s", f.f, f.n, f.f)
		}
	}
}
