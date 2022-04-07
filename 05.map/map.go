package _map

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Map struct {
	mu sync.Mutex
	// 基本上你可以把它看成一个安全的只读的 map
	// 它包含的元素其实也是通过原子操作更新的，但是已删除的 entry 就需要加锁操作了
	read atomic.Value // readOnly

	// 包含需要加锁才能访问的元素
	// 包括所有在 read 字段中但未被 expunged（删除）的元素以及新加的元素
	dirty map[interface{}]*entry

	// 记录从 read 中读取 miss 的次数，一旦 miss 数和 dirty 长度一样了，就会把 dirty 提升为 read，并把 dirty 置空
	misses int
}
type readOnly struct {
	m       map[interface{}]*entry
	amended bool // 当 dirty 中包含 read 没有的数据时为 true，比如新增一条数据
}

// expunged 是用来标识此项已经删掉的指针
// 当 map 中的一个项目被删除了，只是把它的值标记为 expunged，以后才有机会真正删除此项
var expunged = unsafe.Pointer(new(interface{}))

// entry 代表一个值
type entry struct {
	p unsafe.Pointer // *interface{}
}

func (m *Map) Store(key, value interface{}) {
	read, _ := m.read.Load().(readOnly)
	// 如果read字段包含这个项，说明是更新，cas更新项目的值即可
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}
	// read中不存在，或者cas更新失败，就需要加锁访问dirty了
	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok { // 双检查，看看read是否已经存在了
		if e.unexpungeLocked() {
			// 此项目先前已经被删除了，通过将它的值设置为nil，标记为unexpunged
			m.dirty[key] = e
		}
		e.storeLocked(&value) // 更新
	} else if e, ok := m.dirty[key]; ok { // 如果dirty中有此项
		e.storeLocked(&value) // 直接更新
	} else { // 否则就是一个新的key
		if !read.amended { //如果dirty为nil
			// 需要创建dirty对象，并且标记read的amended为true,
			// 说明有元素它不包含而dirty包含
			m.dirtyLocked()
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value) //将新值增加到dirty对象中
	}
	m.mu.Unlock()
}
