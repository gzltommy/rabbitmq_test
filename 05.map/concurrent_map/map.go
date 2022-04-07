package concurrent_map

import "sync"

var SHARD_COUNT = 32

// 分成 SHARD_COUNT 个分片的 map
type ConcurrentMap []*ConcurrentMapShared

// 通过 RWMutex 保护的线程安全的分片，包含一个 map
type ConcurrentMapShared struct {
	items        map[string]interface{}
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// 创建并发 map
func New() ConcurrentMap {
	m := make(ConcurrentMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
	}
	return m
}

// 根据 key 计算分片索引
func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)

	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m ConcurrentMap) Set(key string, value interface{}) {
	// 根据 key 计算出对应的分片
	shard := m.GetShard(key)
	shard.Lock() //对这个分片加锁，执行业务操作
	shard.items[key] = value
	shard.Unlock()
}

func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	// 根据 ke y计算出对应的分片
	shard := m.GetShard(key)
	shard.RLock()
	// 从这个分片读取 key 的值
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}
