package loggo

import (
	"time"
	"sync"
)

type Record struct{
	level Flag 
	gid int64
	time time.Time

	fmt *string
	args []interface{}
}

var (
	record_pool = &sync.Pool{
		New : func() interface{} {
			return &Record{}
		},
	}
)

func GetRecord() *Record {
	return record_pool.Get().(*Record)
}
func PutRecord(r *Record) {
	record_pool.Put(r)
}

