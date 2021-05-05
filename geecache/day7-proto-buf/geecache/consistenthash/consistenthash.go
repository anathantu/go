package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//Hash maps bytes to uint32
type Hash func(data []byte) uint32

//Map contains all hashed keys
type Map struct {
	hash     Hash           //hash函数，利用这个变量可以实现自定义函数
	replicas int            //虚拟节点倍数
	keys     []int          //哈希环
	hashMap  map[int]string //将不同的哈希值与对应节点映射起来
}

//New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，
	//因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
