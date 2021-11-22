package cache

import (
	"hash/fnv"
	"log"
)

type Cache struct {
	cache map[string][]*Sample

	maxBytes int64

	nbytes int64

	n int64
}

// New is the Constructor of Cache
func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		cache:    make(map[string][]*Sample),
	}
}

func (c *Cache) Check(key string, v *Sample) bool {
	num := c.nbytes + int64(len(key)) + 16
	if num > c.maxBytes {
		log.Printf("大小超出限制，现在缓存的数据有%d", c.n)
		return false
	}
	return true
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, v *Sample) {
	if ele, ok := c.cache[key]; ok {
		c.cache[key] = append(ele, v)
	} else {
		nEle := make([]*Sample, 1)
		nEle = append(nEle, v)
		c.cache[key] = nEle
		c.nbytes += int64(len(key))
	}
	c.n++
	c.nbytes += 16
}

type series struct {
	metric_name string

	instance string

	lset []label //其他的label
}

func (s *series) hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(s.metric_name + s.instance))
	for _, l := range s.lset {
		h.Write([]byte(l.name + l.value))
	}
	return h.Sum64()
}

type label struct {
	name  string
	value string
}

//t:timestamp v:value
type Sample struct {
	T int64

	V float64
}
