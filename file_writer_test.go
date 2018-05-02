package loggo

import (
	"os"
	"strings"
	"testing"
	"time"
	"runtime"
)

var (
	test_data   = strings.Repeat("A", 123)
	file_name   = "FileWriter.test"
	file_rename = "FileWriter.test.tmp"
)

func TestWrite(t *testing.T) {
	os.Remove(file_name)
	fw := NewFileWriter(file_name)
	defer os.Remove(fw.file.Name())
	fw.WriteString(test_data)

	time.Sleep(600 * time.Millisecond)
	//fw.file.Sync()
	info, err := fw.file.Stat()
	if err != nil {
		t.Fatalf("file stat fatal. %v", err)
		return
	}
	if info.Size() != 123 {
		t.Fatalf("write file fatal. %d", info.Size())
	}
	fw.Close()
}

func TestClose(t *testing.T) {
	fw := NewFileWriter(file_name)
	defer os.Remove(fw.file.Name())
	fw.WriteString(test_data)

	fw.Close()

	//_, err := fw.file.Stat()
	//if err == nil {
	//	t.Fatalf("file have closed.")
	//	return
	//}

	info, err := os.Stat(fw.file.Name())
	if err != nil {
		t.Fatalf("file stat fatal. %v", err)
		return
	}
	if info.Size() != 123 {
		t.Fatalf("write file fatal. %d", info.Size())
	}
}

func TestReopen(t *testing.T) {
	fw := NewFileWriter(file_name)
	defer os.Remove(fw.file.Name())
	fw.WriteString(test_data)
	os.Rename(file_name, file_rename)
	defer os.Remove(file_rename)
	err := fw.Reopen()
	if err != nil {
		t.Errorf("reopen fail. %v", err)
	}
	t.Logf("reopen %v", err)

	info, err := os.Stat(file_rename)
	if err != nil {
		t.Fatalf("file stat fatal. %v", err)
		return
	}
	if info.Size() != 123 {
		t.Fatalf("write file fatal. %d", info.Size())
	}
	fw.WriteString(test_data)
	fw.Close()

	info, err = os.Stat(fw.file.Name())
	if err != nil {
		t.Fatalf("file stat fatal. %v", err)
		return
	}
	if info.Size() != 123 {
		t.Fatalf("write file fatal. %d", info.Size())
	}
}

func TestReopenFail(t *testing.T) {
	fw := NewFileWriter(file_name)
	defer os.Remove(fw.file.Name())
	fw.WriteString(test_data)

	not_exists_rename := "/tmp/notexist/" + file_rename
	fw.filename = not_exists_rename

	os.Rename(file_name, file_rename)
	defer os.Remove(file_rename)
	err := fw.Reopen()
	if err == nil {
		t.Fatalf("reopen should fail. %v", err)
	}
	t.Logf("reopen %v", err)

	info, err := fw.file.Stat()
	if err != nil {
		t.Fatalf("file stat fatal. %v", err)
		return
	}
	if info.Size() != 123 {
		t.Fatalf("write file fatal. %d", info.Size())
	}

	fw.WriteString(test_data)
	fw.Close()

	info, err = os.Stat(file_rename)
	if err != nil {
		t.Fatalf("file stat fatal. %v", err)
		return
	}
	if info.Size() != 246 {
		t.Fatalf("write file fatal. %d", info.Size())
	}
}

func TestWriteParallel(t *testing.T) {
	os.Remove(file_name)
	fw := NewFileWriter(file_name)
	defer os.Remove(fw.file.Name())

	testdata1 := strings.Repeat("a", 99)
	testdata1 += "\n"
	testdata2 := strings.Repeat("b", 299)
	testdata2 += "\n"
	testdata3 := strings.Repeat("c", 4999)
	testdata3 += "\n"
	testdata4 := strings.Repeat("d", 8999)
	testdata4 += "\n"
	testfunc1 := func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 8; i++ {
			go fw.WriteString(testdata1)
		}
	}
	testfunc2 := func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 4; i++ {
			go fw.WriteString(testdata2)
		}
	}
	testfunc3 := func(t *testing.T) {
		t.Parallel()
		go fw.WriteString(testdata3)
		go fw.WriteString(testdata3)
		go fw.WriteString(testdata3)
	}
	testfunc4 := func(t *testing.T) {
		t.Parallel()
		go fw.WriteString(testdata4)
		go fw.WriteString(testdata4)
	}

	t.Parallel()
	t.Run("group", func(t *testing.T) {
		//t.Parallel()
		t.Run("Test1", testfunc1)
		t.Run("Test2", testfunc2)
		t.Run("Test3", testfunc3)
		t.Run("Test4", testfunc4)
	})

	//time.Sleep(1*time.Second)
	//time.Sleep(100*time.Second)
	runtime.Gosched()
	go fw.WriteString(testdata1)
	fw.Close()
	file, err := os.OpenFile(file_name, os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("OpenFile err. %v", err)
	}
	buf := make([]byte, 35000)
	n, err := file.Read(buf)
	if err != nil {
		t.Fatalf("read err. %v", err)
	}
	if n != 35000 {
		t.Fatalf("read not enough. %d", n)
	}
}
