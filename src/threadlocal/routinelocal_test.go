package threadlocal

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

const LOOP_CNT = 1000

func TestGet(t *testing.T) {
	var waitGroup sync.WaitGroup
	for i := 0; i < LOOP_CNT; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			Set("probe", i)
			defer Clear()
			assert.Equal(t, Get("probe"), i)
		}(i)
	}
	waitGroup.Wait()
}

func TestGetWrong(t *testing.T) {
	var waitGroup sync.WaitGroup
	for i := 0; i < LOOP_CNT; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			Set("probe", i)
			defer Clear()
			assert.NotEqual(t, Get("wrong-key"), i)
		}(i)
	}
	waitGroup.Wait()
}

func TestRoutineLocal_Remove(t *testing.T) {
	var waitGroup sync.WaitGroup
	for i := 0; i < LOOP_CNT; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			Set("probe", i)
			defer Clear()
			assert.Equal(t, Get("probe"), i)
			Remove("probe")
			assert.Nil(t, Get("probe"))
		}(i)
	}
	waitGroup.Wait()
}

func TestGetCurGoroutineID(t *testing.T) {
	routineID := GetCurGoroutineID()
	assert.NotEmpty(t, routineID)
	go func() {
		assert.NotEmpty(t, GetCurGoroutineID())
	}()
	println(routineID)
}

func BenchmarkGetCurGoroutineID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCurGoroutineID()
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var waitGroup sync.WaitGroup
		for i := 0; i < LOOP_CNT; i++ {
			waitGroup.Add(1)
			go func(i int) {
				defer waitGroup.Done()
				Set("probe", i)
				defer Clear()
				Get("probe")
			}(i)
		}
		waitGroup.Wait()
	}
}
