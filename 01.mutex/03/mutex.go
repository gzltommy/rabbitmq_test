package mutex

import (
	"runtime"
	"sync/atomic"
)

type Mutex struct {
	state int32
	sema  uint32
}

const (
	mutexLocked      = 1 << iota // mutex is locked：1
	mutexWoken                   // 2
	mutexWaiterShift = iota      // 2
)

func (m *Mutex) Lock() {
	// Fast path: 幸运之路，正好获取到锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	awoke := false
	iter := 0
	for { // 不管是新来的请求锁的goroutine, 还是被唤醒的goroutine，都不断尝试请求锁
		old := m.state            // 先保存当前锁的状态
		new := old | mutexLocked  // 新状态设置加锁标志
		if old&mutexLocked != 0 { // 锁还没被释放
			if runtime_canSpin(iter) { // 还可以自旋
				if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
					atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
					awoke = true
				}
				runtime_doSpin()
				iter++
				continue // 自旋，再次尝试请求锁
			}
			new = old + 1<<mutexWaiterShift
		}
		if awoke { // 唤醒状态
			if new&mutexWoken == 0 {
				panic("sync: inconsistent mutex state")
			}
			new &^= mutexWoken // 新状态清除唤醒标记
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&mutexLocked == 0 { // 旧状态锁已释放，新状态成功持有了锁，直接返回
				break
			}
			runtime_Semacquire(&m.sema) // 阻塞等待
			awoke = true                // 被唤醒
			iter = 0
		}
	}
}

func (m *Mutex) Unlock() {
	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked) //去掉锁标志
	if (new+mutexLocked)&mutexLocked == 0 {        //本来就没有加锁
		panic("sync: unlock of unlocked mutex")
	}

	old := new
	for {
		if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 { // 没有等待者，或者有唤醒的waiter，或者锁原来已加锁
			return
		}
		new = (old - 1<<mutexWaiterShift) | mutexWoken // 新状态，准备唤醒goroutine，并设置唤醒标志
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			runtime.Semrelease(&m.sema)
			return
		}
		old = m.state
	}
}
