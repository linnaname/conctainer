package threadsafe

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	ROUTINE_NUM = 1000
	INNER_LOOP  = 5237659
)

func TestLongAdder_Sum(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		routineNum int
		innerLop   int
		expected   int
	}{
		//{2000, 5237659, 2000 * 5237659},
		{10, 100, 10 * 100},
		//{1, 1, 1 * 1},
		//{5, 3, 5 * 3},
		//{999, 888, 999 * 888},
	}

	for _, test := range tests {
		adder := New()
		var wg sync.WaitGroup

		for i := 0; i < test.routineNum; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < test.innerLop; j++ {
					adder.Increment()
				}
				wg.Done()
			}()
		}
		wg.Wait()

		assert.NotEqual(adder.Sum(), 0)
		assert.Equal(adder.Sum(), int64(test.expected))
	}

}

func TestAtomic_Sum(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		routineNum int
		innerLop   int
		expected   int
	}{
		//{2000, 5237659, 2000 * 5237659},
		{500, 655351, 500 * 655351},
		//{1, 1, 1 * 1},
		//{5, 3, 5 * 3},
		//{999, 888, 999 * 888},
	}

	for _, test := range tests {
		var wg sync.WaitGroup
		val := int64(0)
		for i := 0; i < test.routineNum; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < test.innerLop; j++ {
					atomic.AddInt64(&val, 1)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		assert.NotEqual(atomic.LoadInt64(&val), 0)
		assert.Equal(atomic.LoadInt64(&val), int64(test.expected))

	}

}

func BenchmarkLongAdder_Sum(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var tests = []struct {
			routineNum int
			innerLop   int
			expected   int
		}{
			//{2000, 5237659, 2000 * 5237659},
			{200, 200, 200 * 200},
			//{1, 1, 1 * 1},
			//{5, 3, 5 * 3},
			//{999, 888, 999 * 888},
		}

		for _, test := range tests {
			adder := New()
			var wg sync.WaitGroup

			for i := 0; i < test.routineNum; i++ {
				wg.Add(1)
				go func() {
					for j := 0; j < test.innerLop; j++ {
						adder.Increment()
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}
}

func BenchmarkAtomic_Sum(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var tests = []struct {
			routineNum int
			innerLop   int
			expected   int
		}{
			//{2000, 5237659, 2000 * 5237659},
			{200, 200, 200 * 200},
			//{1, 1, 1 * 1},
			//{5, 3, 5 * 3},
			//{999, 888, 999 * 888},
		}

		for _, test := range tests {
			var wg sync.WaitGroup
			val := int64(0)
			for i := 0; i < test.routineNum; i++ {
				wg.Add(1)
				go func() {
					for j := 0; j < test.innerLop; j++ {
						atomic.AddInt64(&val, 1)
					}
					wg.Done()
				}()
			}
			wg.Wait()
			atomic.LoadInt64(&val)
		}
	}

}
