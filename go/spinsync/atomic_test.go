// Source code file, created by Developer@YAN_YING_SONG.

// go test -bench=Benchmark*

package spinsync

import (
	"runtime"
	"sync"
	"testing"
)

var _ = func() struct{} {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return struct{}{}
}()

func Benchmark_syncMap(b *testing.B) {
	var wg = &sync.WaitGroup{}
	sets := &sync.Map{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			sets.Store("00000000", "11111111")
			sets.Delete("00000000")
		}(wg)
	}
	wg.Wait()
}

func Benchmark_mutexMap(b *testing.B) {
	var wg = &sync.WaitGroup{}
	sets := make(map[string]string, 10)
	sLock := &sync.Mutex{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			sLock.Lock()
			sets["00000000"] = "11111111"
			sLock.Unlock()
			sLock.Lock()
			delete(sets, "00000000")
			sLock.Unlock()
		}(wg)
	}
	wg.Wait()
}

func Benchmark_spinMap(b *testing.B) {
	var wg = &sync.WaitGroup{}
	sets := make(map[string]string, 10)
	sLock := &Spinlock{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			sLock.Lock()
			sets["00000000"] = "11111111"
			sLock.Unlock()
			sLock.Lock()
			delete(sets, "00000000")
			sLock.Unlock()
		}(wg)
	}
}

/*

PS C:\Work\library\generic\spinsync> go test -bench=Benchmark*
goos: windows
goarch: amd64
pkg: library/generic/spinsync
cpu: AMD Ryzen 7 6800HS with Radeon Graphics
Benchmark_syncMap-16              583010              2230 ns/op
Benchmark_mutexMap-16            1000000              1010 ns/op
Benchmark_spinMap-16             2408223               560.4 ns/op
PASS
ok      library/generic/spinsync        4.656s

Mac-mini % go test -bench='Benchmark*'
goos: darwin
goarch: arm64
pkg: libagent/generic/spinsync
cpu: Apple M2
Benchmark_syncMap-8       906150              1116 ns/op
Benchmark_mutexMap-8     3644122               341.3 ns/op
Benchmark_spinMap-8      4318791               281.5 ns/op
PASS
ok      libagent/generic/spinsync       5.831s

*/
