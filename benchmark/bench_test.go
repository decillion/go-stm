package bench

import (
	"sync"
	"sync/atomic"
	"testing"

	de "github.com/decillion/go-stm"
	lu "github.com/lukechampine/stm"
)

func Benchmark_Read100_RWMutex(b *testing.B) {
	var x, y int
	mu := &sync.RWMutex{}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = x
			_ = y
			mu.RUnlock()
		}
	})
}
func Benchmark_Read100_decillion(b *testing.B) {
	x := de.New(0)
	y := de.New(0)

	load := func(rec *de.TRec) interface{} {
		rec.Load(x)
		rec.Load(y)
		return nil
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			de.Atomically(load)
		}
	})
}

func Benchmark_Read100_lukechampine(b *testing.B) {
	x := lu.NewVar(0)
	y := lu.NewVar(0)

	load := func(tx *lu.Tx) {
		tx.Get(x)
		tx.Get(y)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lu.Atomically(load)
		}
	})
}

func Benchmark_Read90Write10_RWMutex(b *testing.B) {
	var x, y int
	mu := &sync.RWMutex{}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				mu.Lock()
				x++
				y++
				mu.Unlock()
			} else {
				mu.RLock()
				_ = x
				_ = y
				mu.RUnlock()
			}
		}
	})
}

func Benchmark_Read90Write10_decillion(b *testing.B) {
	x := de.New(0)
	y := de.New(0)

	inc := func(rec *de.TRec) interface{} {
		rec.Store(x, rec.Load(x).(int)+1)
		rec.Store(y, rec.Load(y).(int)+1)
		return nil
	}
	load := func(rec *de.TRec) interface{} {
		rec.Load(x)
		rec.Load(y)
		return nil
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				de.Atomically(inc)
			} else {
				de.Atomically(load)
			}
		}
	})
}

func Benchmark_Read90Write10_lukechampine(b *testing.B) {
	x := lu.NewVar(0)
	y := lu.NewVar(0)

	inc := func(tx *lu.Tx) {
		tx.Set(x, tx.Get(x).(int)+1)
		tx.Set(y, tx.Get(y).(int)+1)
	}
	load := func(tx *lu.Tx) {
		tx.Get(x)
		tx.Get(y)
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				lu.Atomically(inc)
			} else {
				lu.Atomically(load)
			}
		}
	})
}
