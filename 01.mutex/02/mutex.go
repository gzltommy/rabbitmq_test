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
	// Fast path: 幸运 case，能够直接获取到锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	awoke := false
	for {
		old := m.state
		new := old | mutexLocked // 新状态加锁（如果能成功地通过 CAS 把这个新值赋予 state（第 38 行和第 39 行），就代表抢夺锁的操作成功了）
		if old&mutexLocked != 0 {
			new = old + 1<<mutexWaiterShift //等待者数量加一
		}
		if awoke {
			// goroutine是被唤醒的，
			// 新状态清除唤醒标志
			new &^= mutexWoken
		}

		//
		if atomic.CompareAndSwapInt32(&m.state, old, new) { //设置新状态
			if old&mutexLocked == 0 { // 锁原状态未加锁
				break
			}
			runtime.Semacquire(&m.sema) // 请求信号量,休眠
			awoke = true                // 休眠醒来之后 awoke 置为 true，尝试争抢锁。
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
