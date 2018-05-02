package loggo

import (

)

type MultiBackend struct {
	//backends []Logger
}

func (mb *MultiBackend) log(lv Flag, format *string, args ...interface{}) {
	//for _, b := range mb.backends {
	//	b.log(lv, 1, format, args)
	//}
}
