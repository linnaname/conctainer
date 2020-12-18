package threadsafe

import (
	"reflect"
	"runtime"
	"sync/atomic"
)

type Striped64 struct {
	base      int64
	cellsBusy int32
	cells     Cells
}

func (s *Striped64) casBase(cmp, val int64) bool {
	return atomic.CompareAndSwapInt64(&s.base, cmp, val)
}

func (s *Striped64) casCellsBusy() bool {
	return atomic.CompareAndSwapInt32(&s.cellsBusy, 0, 1)
}

func (s *Striped64) longAccumulate(probe int, val int64, uncontended bool) {
	h := GetProbe()
	if h == 0 {
		h = GetProbe()
		uncontended = true
	}

	collide := false
	for {
		as := s.cells
		n := len(as)
		if as != nil && n > 0 {
			a := as[(n-1)&h]
			if a == nil {
				if atomic.LoadInt32(&s.cellsBusy) == 0 {
					r := &Cell{val: val}
					if atomic.LoadInt32(&s.cellsBusy) == 0 && s.casCellsBusy() {
						var created bool = false
						rs := s.cells
						m := len(rs)
						if rs != nil && m > 0 {
							j := (m - 1) & h
							if rs[j] == nil {
								rs[j] = r
								created = true
							}
						}
						atomic.StoreInt32(&s.cellsBusy, 0)
						if created {
							break
						}
						continue
					}
				}
				collide = false
			} else if !uncontended {
				uncontended = true
			} else if a.cas(a.val, a.val+val) {
				break
			} else if n >= runtime.NumCPU() || !reflect.DeepEqual(s.cells, as) {
				collide = false
			} else if !collide {
				collide = true
			} else if atomic.LoadInt32(&s.cellsBusy) == 0 && s.casCellsBusy() {
				if reflect.DeepEqual(s.cells, as) { // Expand table unless stale
					rs := make(Cells, n<<1, n<<2)
					copy(rs, as)
					s.cells = rs
				}
				atomic.StoreInt32(&s.cellsBusy, 0)
				collide = false
				continue
			}
			h = AdvanceProbe(h)
		} else if atomic.LoadInt32(&s.cellsBusy) == 0 && reflect.DeepEqual(s.cells, as) && s.casCellsBusy() {
			init := false
			if reflect.DeepEqual(s.cells, as) {
				rs := make(Cells, 2)
				rs[h&1] = &Cell{val: val}
				s.cells = rs
				init = true
			}
			atomic.StoreInt32(&s.cellsBusy, 0)
			if init {
				break
			}
		} else if s.casBase(s.base, s.base+val) {
			break
		}
	}
}
