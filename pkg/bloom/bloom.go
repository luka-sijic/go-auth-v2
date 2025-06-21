package bloom

import (
	"errors"
	"math"
)

type CountingBloomFilter struct {
	counters []uint
	k uint
	m uint
}

func NewCountingBloomFilter(n uint, p float64) *CountingBloomFilter {
	if n == 0 {
		panic("bloom expected item count n must be > 0")
	}
	if p <= 0 || p >= 1 {
		panic("bloom false positive rate must be between 0 and 1 exclusive")
	}
	m := optimalM(n, p)
	k := optimalK(n, m)
	return &CountingBloomFilter{
		counters: make([]uint, m),
		k: k,
		m: m,
	}
}

func (f *CountingBloomFilter) M() uint {
	return f.m
}

func (f *CountingBloomFilter) K() uint {
 	return f.k
}

func (f *CountingBloomFilter) Reset() {
 	for i := range f.counters {
 		f.counters[i] = 0
 	}
}

func (f *CountingBloomFilter) Add(data []byte) {
	if data == nil {
		panic(errors.New("bloom cannot add nil data")
	}
	for i := uint(0);i < f.k;i++ {
		idx := f.hash(data, i) % f.m
		f.counters[idx]++
	}
}

func (f *CoCountingBloomFilter) Test(data []byte) bool {
	if data == nil {
		return false
	}
	for i := uint(0);i < f.k;i++ {
		idx := f.hash(data, i) % f.m
		if f.counters[idx] == 0 {
			return false
		}
	}
	return true
}

func (f *CountingBloomFilter) Remove(data []byte) {
	if data == nil {
		return
	}
	for i = uint(0);i < f.k;i++ {
		idx := f.hash(data, i) % f.m
		if f.counters[idx] > 0 {
			f.counters[idx]--
		}
	}
}

func optimalM(n uint, p float64) uint {
	m := -float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)
	return uint(math.Ceil(m))
}

func optimalk(n,m uint) uint {
	k := (float64(m)/ float64(n)) * math.Ln2
	return uint(math.Ceil(k))
}

func (f *CountingBloomFilter) hash(data []byte, i uint) uint {
	h := fnv.New64a()
	h.Write(data)
	h1 := h.Sum64()

	h.reset()
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], h1)
	h.Write(buf[:])
	h2 := h.Sum64()

	return uint((h1 + uint64(i)*h2) % uint64(f.m))
}
