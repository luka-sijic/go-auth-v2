package bloom

type CBF struct {
	count     []uint32
	size      uint32
	hashCount uint32
}

func (c *CBF) Init(s, h uint32) {
	c.count = make([]uint32, int(s))
	c.size = s
	c.hashCount = h
}

func (c *CBF) Insert(data string) {
	h1 := murmurHash(data, 0x9747b28c)
	h2 := murmurHash(data, 0x12345678)

	for i := 0; i < int(c.hashCount); i++ {
		index := c.hash(data, uint32(i), h1, h2)
		c.count[index]++
	}
}

func (c *CBF) PossiblyContains(data string) bool {
	h1 := murmurHash(data, 0x9747b28c)
	h2 := murmurHash(data, 0x12345678)

	for i := 0; i < int(c.hashCount); i++ {
		index := c.hash(data, uint32(i), h1, h2)
		if c.count[index] == 0 {
			return false
		}
	}
	return true
}

func (c *CBF) hash(key string, i, h1, h2 uint32) uint32 {
	return (h1 + i*h2) % c.size
}

func murmurHash(key string, seed uint32) uint32 {
	const m uint32 = 0x5bd1e995 // magic constant from Murmur2
	const r uint32 = 24

	h := seed ^ uint32(len(key))

	data := []byte(key)
	for len(data) >= 4 {
		k := uint32(data[0]) |
			uint32(data[1])<<8 |
			uint32(data[2])<<16 |
			uint32(data[3])<<24

		k *= m
		k ^= k >> r
		k *= m

		h *= m
		h ^= k

		data = data[4:]
	}

	switch len(data) {
	case 3:
		h ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint32(data[0])
		h *= m
	}

	// Final avalanche.
	h ^= h >> 13
	h *= m
	h ^= h >> 15

	return h
}
