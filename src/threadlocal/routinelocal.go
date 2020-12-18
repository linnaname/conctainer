package threadlocal

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

const BYTE_SIZE = 64

var defaultLocal = NewRoutineLocal()
var goroutineSpace = []byte("goroutine ")

func Get(k string) interface{} {
	return defaultLocal.Get(k)
}

func Set(k string, v interface{}) {
	defaultLocal.Set(k, v)
}

func Remove(k string) {
	defaultLocal.Remove(k)
}

func Clear() {
	defaultLocal.Clear()
}

/**
参考golang的源代码
https://github.com/golang/net/blob/master/http2/gotrack.go#L51
*/
func GetCurGoroutineID() uint64 {
	b := make([]byte, BYTE_SIZE)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

type IRoutineLocal interface {
	Set(string, interface{})
	Get(string) interface{}
	Remove(string)
	Clear()
}

type RoutineLocal struct {
	val sync.Map
}

func NewRoutineLocal() IRoutineLocal {
	return &RoutineLocal{val: sync.Map{}}
}

func (r *RoutineLocal) Set(k string, v interface{}) {
	routineID := GetCurGoroutineID()
	vmap, ok := r.val.Load(routineID)
	if !ok {
		vmap = make(map[string]interface{})
		r.val.Store(routineID, vmap)
	}
	vmap.(map[string]interface{})[k] = v
}

func (r *RoutineLocal) Get(k string) interface{} {
	routineID := GetCurGoroutineID()
	vmap, ok := r.val.Load(routineID)
	if !ok {
		return nil
	}
	return vmap.(map[string]interface{})[k]
}

func (r *RoutineLocal) Remove(k string) {
	routineID := GetCurGoroutineID()
	vm, ok := r.val.Load(routineID)
	if !ok {
		return
	}
	delete(vm.(map[string]interface{}), k)
}

func (r *RoutineLocal) Clear() {
	r.val.Delete(GetCurGoroutineID())
}
