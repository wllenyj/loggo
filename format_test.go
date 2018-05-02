package loggo

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func test(format Formatter, record *Record) string {
	return test1(format, record)
}
func test1(format Formatter, record *Record) string {
	return test2(format, record)
}
func test2(format Formatter, record *Record) string {
	return test3(format, record)
}
func test3(format Formatter, record *Record) string {
	return test4(format, record)
}
func test4(format Formatter, record *Record) string {
	return test5(format, record)
}
func test5(format Formatter, record *Record) string {
	return test6(format, record)
}
func test6(format Formatter, record *Record) string {
	return FormatterProxy(format, 0, record)
}

var (
	aaa = "Hello,%s!"
	testTime, _ = time.Parse(time.RFC3339Nano, "2018-04-27T19:22:01.40731219+08:00")
	testRecord  = []*Record{
		&Record{DEBUG, 123, testTime, &aaa, []interface{}{"world"}},
		&Record{WARNING, 23652523, testTime, nil, []interface{}{"Hello,world!"}},
	}
	testFormat = []string{
		"%{message}",
		"[%{level}] %{callpath:3} ===",
		"%{time:01-02T15:04:05.999999} %{level} %{shortfile} %{shortfunc} %{message}",
		"%{%{}}",
		"%{unkown}",
	}

	formater = MustStringFormatter(testFormat[2])
)

func TestFormater(t *testing.T) {
	t.Parallel()
	result := test(formater, testRecord[0])
	if result != "04-27T19:22:01.407312 DEBUG format_test.go:30 test6 Hello,world!\n" {
		t.Fatalf("failt: %s", result)
	}
	t.Logf(result)
}

func TestFormatter(t *testing.T) {
	//tests := []struct{
	//	record *TestRecord
	//	format string
	//	result string
	//}
}

func TestNewFormatter(t *testing.T) {
	//"%{message}",
	//"[%{level}] %{callpath:3} ===",
	//"%{time:01-02T15:04:05.999999} %{level} %{shortfile} %{shortfunc} %{message}",
	//"%{%{}}",
	testData := []struct {
		format string
		parts  []part
		err    string
	}{
		{testFormat[0], []part{{fmtVerbMessage, "%s"}}, ""},
		{testFormat[1], []part{
			{fmtVerbStatic, "["},
			{fmtVerbLevel, "%s"},
			{fmtVerbStatic, "] "},
			{fmtVerbCallpath, "3"},
			{fmtVerbStatic, " ==="},
		}, ""},
		{testFormat[3], nil, "invalid log format"},
		{testFormat[4], nil, "unknown variable"},
	}
	for _, test := range testData {
		formatter, err := NewStringFormatter(test.format)
		if test.err != "" {
			if !strings.Contains(err.Error(), test.err) {
				t.Fatalf("New fail. %v %v", err, test.err)
			}
			continue
		} else if err != nil {
			t.Fatalf("New fail. %v", err)
		}
		f := formatter.(*stringFormatter)
		if len(f.parts) != len(test.parts) {
			t.Fatalf("%s fatal. len %d != %d", test.format, len(f.parts), len(test.parts))
		}
		for i := 0; i < len(f.parts); i++ {
			if f.parts[i].verb != test.parts[i].verb {
				t.Fatalf("%s verb fatal. %#v", test.format, f.parts[i].verb)
			}
			if f.parts[i].layout != test.parts[i].layout {
				t.Fatalf("%s layout fatal. %#v", test.format, f.parts[i].layout)
			}
		}

	}
}

func TestParallel(t *testing.T) {
	testfunc := func(t *testing.T) {
		t.Parallel()
		result := test(formater, testRecord[0])
		if strings.Compare(result, "04-27T19:22:01.407312 DEBUG format_test.go:30 test6 Hello,world!\n") != 0{
			t.Fatalf("failt: %x", result)
		}
	}
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", testfunc)
		t.Run("Test2", testfunc)
		t.Run("Test3", testfunc)
		t.Run("Test4", testfunc)
	})
}

func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = test(formater, testRecord[0])
		}
	})
}

func BenchmarkFormatStringsBuilder(b *testing.B) {
	//formater := MustStringFormatter(testFormat[2])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		formater.Format(0, testRecord[0], &output)
		output.String()
	}
}
func BenchmarkFormatByetsBuffer(b *testing.B) {
	//formater := MustStringFormatter(testFormat[2])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output bytes.Buffer
		formater.Format(0, testRecord[0], &output)
		output.String()
	}
}

var (
	builder_pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	buffer_pool_ = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

func BenchmarkFormatStringsBuilderPoolNoDefer(b *testing.B) {
	//formater := MustStringFormatter(testFormat[2])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := builder_pool.Get().(*strings.Builder)
		output.Reset()
		formater.Format(0, testRecord[0], output)
		output.String()
		builder_pool.Put(output)
	}
}
func BenchmarkFormatByetsBufferPoolNoDefer(b *testing.B) {
	//formater := MustStringFormatter(testFormat[2])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := buffer_pool_.Get().(*bytes.Buffer)
		output.Reset()
		formater.Format(0, testRecord[0], output)
		output.String()
		buffer_pool_.Put(output)
	}
}
func BenchmarkFormatStringsBuilderPool(b *testing.B) {
	//formater := MustStringFormatter(testFormat[2])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := builder_pool.Get().(*strings.Builder)
		defer builder_pool.Put(output)
		output.Reset()
		formater.Format(0, testRecord[0], output)
		output.String()
	}
}
func BenchmarkFormatByetsBufferPool(b *testing.B) {
	//formater := MustStringFormatter(testFormat[2])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := buffer_pool_.Get().(*bytes.Buffer)
		defer buffer_pool_.Put(output)
		output.Reset()
		formater.Format(0, testRecord[0], output)
		output.String()
	}
}
