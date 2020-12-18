package threadsafe

import (
	"sync/atomic"
)

type LongAdder struct {
	Striped64
}

func New() *LongAdder {
	return &LongAdder{}
}

func (l *LongAdder) Add(val int64) {
	as := l.cells
	b := atomic.LoadInt64(&l.base)
	if as != nil || !l.casBase(b, b+val) {
		uncontended := true
		if as == nil {
			l.longAccumulate(GetProbe(), val, uncontended)
			return
		}

		m := len(as) - 1
		if m < 0 {
			l.longAccumulate(GetProbe(), val, uncontended)
			return
		}

		probe := GetProbe() & m
		a := l.cells[probe]
		if a == nil {
			l.longAccumulate(probe, val, uncontended)
		} else {
			v := atomic.LoadInt64(&a.val)
			if uncontended = a.cas(v, v+val); !uncontended {
				l.longAccumulate(probe, val, uncontended)
			}
		}
	}
}

func (l *LongAdder) Increment() {
	l.Add(1)
}

func (l *LongAdder) Decrement() {
	l.Add(-1)
}

func (l *LongAdder) Sum() int64 {
	as := l.cells
	sum := l.base
	if as != nil {
		for i := 0; i < len(as); i++ {
			a := as[i]
			if a != nil {
				sum += a.val
			}
		}
	}
	return sum
}

func (l *LongAdder) Reset() {
	as := l.cells
	atomic.StoreInt64(&l.base, 0)
	if as != nil {
		for i := 0; i < len(as); i++ {
			a := as[i]
			if a != nil {
				a.val = 0
			}
		}
	}
}

func (l *LongAdder) SumThenReset() int64 {
	as := l.cells
	sum := atomic.LoadInt64(&l.base)
	atomic.StoreInt64(&l.base, 0)
	if as != nil {
		for i := 0; i < len(as); i++ {
			a := as[i]
			if a != nil {
				sum += a.val
				a.val = 0
			}
		}
	}
	return sum
}

func (l *LongAdder) Int64() int64 {
	return l.Sum()
}
