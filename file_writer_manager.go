package loggo

import (
	"path/filepath"
	"sync"
	"fmt"
)

var (
	manager *FileWriterManager
)

type FileWriterManager struct {
	m *sync.Map
}

func GetFileWriter(filename string) BufferWriter {
	path, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	return manager.getFileWriter(path)
}

func (manager *FileWriterManager) getFileWriter(path string) *FileWriter {
	value, ok := manager.m.Load(path)
	var fw *FileWriter
	if !ok {
		fw = newFileWriter(path, buffer_pool)
		manager.m.Store(path, fw)
	} else {
		fw = value.(*FileWriter)
	}
	return fw 
}

func Reopen() error {
	return manager.reopen()
}

func (manager *FileWriterManager) reopen() error {
	var reterr error
	manager.m.Range(func(key, value interface{}) bool {
		fw := value.(*FileWriter)
		err := fw.Reopen()
		if err != nil {
			reterr = fmt.Errorf("%s reopen err: %v. \n", key.(string), err)	
		}
		return true
	})
	return reterr
}

func Close() {
	manager.Close()
}

func (manager *FileWriterManager) Close() {
	manager.m.Range(func(key, value interface{}) bool {
		fw := value.(*FileWriter)
		fw.Close()
		//if err != nil {
		//	reterr = fmt.Errorf("%s reopen err: %v. \n", key.(string), err)	
		//}
		return true
	})
}

func init() {
	manager = &FileWriterManager{
		m: &sync.Map{},
	}
}
