package threadsafe

import (
	"math/rand"
	"sync"
	"time"
)

type Probe struct {
	val int
}

var pool sync.Pool

func GetProbe() int {
	v := pool.Get()
	if v == nil {
		v = &Probe{}
	}
	r := v.(*Probe)
	x := r.Probe()
	pool.Put(r)
	return x
}

func (p *Probe) Probe() int {
	if p.val == 0 {
		p.initProbe()
	}
	val := AdvanceProbe(p.val)
	p.val = val
	return val
}

func (p *Probe) initProbe() {
	rand.Seed(time.Now().UnixNano())
	val := rand.Int()
	//skip 0
	if p.val == 0 {
		p.val = 1
	}
	p.val = val
}

func AdvanceProbe(probe int) int {
	probe ^= probe << 13 // xorshift
	probe ^= probe >> 17
	probe ^= probe << 5
	return probe
}
