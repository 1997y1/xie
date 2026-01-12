// Source code file, created by Developer@YAN_YING_SONG.

// Faster than sync.Mutex!

package spinsync

import (
	"runtime"
	"sync/atomic"
)

// A more appropriate value obtained by testing.
// const maxBackoff = 3

type Spinlock struct {
	state uint32
}

func (sl *Spinlock) Lock() {
	// If the lock is already in use, the spin remains waiting.

	// func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
	//
	// if *addr == old {
	//     *addr = new
	//     return true
	// }
	//
	// return false // *addr != old

	// backoff := 0
	for !atomic.CompareAndSwapUint32(&sl.state, 0, 1) {
		runtime.Gosched()
		// AutoSchedule(&backoff)
	}
}

func (sl *Spinlock) Unlock() {
	// Release the lock and give it to other threads.

	atomic.StoreUint32(&sl.state, 0)
}

func (sl *Spinlock) Atomic(f func()) {
	// Integrate the syntactic sugar of Lock Unlock.

	sl.Lock()
	defer sl.Unlock()
	f()
}

// func AutoSchedule(backoff *int) {
// 	// Automatic scheduling to avoid too frequent runtime.Gosched().
//
// 	*backoff++
// 	if *backoff >= maxBackoff {
// 		*backoff = 0
// 		runtime.Gosched()
// 	}
// }
