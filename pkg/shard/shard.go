package shard

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashRing struct {
	replicas int
	keys     []int
	hashMap  map[int]string
	mutex    sync.RWMutex
}

func NewHashRing(replicas int, servers []string) *HashRing {
	hr := &HashRing{
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	for _, server := range servers {
		hr.Add(server)
	}
	return hr
}

func (hr *HashRing) Add(server string) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	for i := 0; i < hr.replicas; i++ {
		hash := int(crc32.ChecksumIEEE([]byte(server + strconv.Itoa(i))))
		hr.keys = append(hr.keys, hash)
		hr.hashMap[hash] = server
	}
	sort.Ints(hr.keys)
}

func (hr *HashRing) Get(key string) string {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	if len(hr.keys) == 0 {
		return ""
	}
	hash := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.Search(len(hr.keys), func(i int) bool { return hr.keys[i] >= hash })

	if idx == len(hr.keys) {
		idx = 0
	}
	return hr.hashMap[hr.keys[idx]]
}
