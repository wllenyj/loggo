package loggo

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	DEF_CHAN_LEN             = 1024
	DEF_FILE_PEM os.FileMode = 0644
	DEF_BUF_SIZE             = 8 * 1024
)

type FileWriter struct {
	filename string
	bw       *bufio.Writer
	file     *os.File
	closed   uint32

	ch     chan []byte
	quit   chan chan error
	reopen chan chan error
}

func (fw *FileWriter) Write(b []byte) (int, error) {
	if atomic.LoadUint32(&fw.closed) == 1 {
		return 0, fmt.Errorf("closed")
	}
	fw.ch <- b
	return len(b), nil
}
func (fw *FileWriter) WriteString(b string) (int, error) {
	if atomic.LoadUint32(&fw.closed) == 1 {
		//fmt.Printf("use closed chan.\n")
		return 0, fmt.Errorf("closed")
	}
	fw.ch <- []byte(b)
	return len(b), nil
}

func (fw *FileWriter) Close() {
	//close(fw.ch)
	//<-fw.quit
	//fw.file.Close()

	atomic.StoreUint32(&fw.closed, 1)
	quited := make(chan error)
	fw.quit <- quited
	err := <-quited
	if err != nil {
		fmt.Printf("LOGGO %s close ret: %v\n", fw.filename, err)
	}
}

func (fw *FileWriter) Reopen() error {
	finish := make(chan error)
	fw.reopen <- finish
	err := <-finish
	return err
}

func NewFileWriter(filename string) *FileWriter {
	bw, file, err := newBufWriter(filename)
	if err != nil {
		panic(fmt.Sprintf("newBufWriter err. %v", err))
	}
	fw := &FileWriter{
		filename: filename,
		file:     file,
		bw:       bw,
		closed:   0,
		ch:       make(chan []byte, DEF_CHAN_LEN),
		//quit:     make(chan struct{}),
		quit:   make(chan chan error),
		reopen: make(chan chan error),
	}
	syn := make(chan struct{})
	go fw.loop(syn)
	<-syn
	return fw
}

func newBufWriter(filename string) (*bufio.Writer, *os.File, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, DEF_FILE_PEM)
	if err != nil {
		return nil, nil, err
	}
	w := bufio.NewWriterSize(file, DEF_BUF_SIZE)
	return w, file, nil
}

func (fw *FileWriter) loop(syn chan struct{}) {
	flushTimer := time.NewTicker(time.Millisecond * 500)
	close(syn)
	for {
		//bw := fw.bw
		select {
		case d := <-fw.ch:
			//case d, ok := <-fw.ch:
			//if !ok {
			//	goto END_FOR
			//}
			if _, err := fw.bw.Write(d); err != nil {
				fmt.Printf("LOGGO ERROR log file write err. %v\n", err)
			}
		case quited := <-fw.quit:
			var reterr error
			lost := len(fw.ch)
			//fmt.Printf("quit lost %d\n", lost)
			for i := 0; i < lost; i++ {
				d := <-fw.ch
				fw.bw.Write(d)
			}
			if err := fw.bw.Flush(); err != nil {
				reterr = fmt.Errorf("LOGGO ERROR log file flush err. %v", err)
			}
			quited <- reterr
		case <-flushTimer.C:
			if err := fw.bw.Flush(); err != nil {
				fmt.Printf("LOGGO ERROR log file flush err. %v\n", err)
			}
		case finish := <-fw.reopen:
			bw, file, reterr := fw.doReopen()
			if bw != nil && file != nil {
				fw.bw = bw
				fw.file = file
			}
			finish <- reterr
		}
	}
	//END_FOR:
	//if err := fw.bw.Flush(); err != nil {
	//	fmt.Printf("log quit flush err. %v\n", err)
	//}
	//close(fw.quit)
}

func (fw *FileWriter) doReopen() (*bufio.Writer, *os.File, error) {
	var reterr error
	if err := fw.bw.Flush(); err != nil {
		reterr = fmt.Errorf("LOGGO ERROR reopen flush err: %v.", err)
	}
	bw, file, err := newBufWriter(fw.filename)
	if err != nil {
		reterr = fmt.Errorf("LOGGO ERROR log reopen newbuf err: %v. %v", err, reterr)
		return nil, nil, reterr
	}
	if err := fw.file.Close(); err != nil {
		reterr = fmt.Errorf("LOGGO ERROR log reopen close err: %v. %v", err, reterr)
	}
	return bw, file, reterr
}
